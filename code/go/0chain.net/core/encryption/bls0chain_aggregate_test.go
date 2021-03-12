package encryption

import (
	"fmt"
	"github.com/herumi/bls/ffi/go/bls"
	"reflect"
	"sync"
	"testing"
)

func TestAggregateSignatures(t *testing.T) {
	total := 1000
	batchSize := 250
	numBatches := total / batchSize
	sigSchemes := make([]SignatureScheme, total)
	msgs := make([]string, total)
	msgHashes := make([]string, total)
	msgSignatures := make([]string, total)
	clientSignatureScheme := "bls0chain"
	for i := 0; i < total; i++ {
		sigSchemes[i] = GetSignatureScheme(clientSignatureScheme)
		sigSchemes[i].GenerateKeys()
		msgs[i] = fmt.Sprintf("testing aggregate messages : %v", i)
		msgHashes[i] = Hash(msgs[i])
		sig, err := sigSchemes[i].Sign(msgHashes[i])
		if err != nil {
			t.Fatal(err)
		}
		msgSignatures[i] = sig
	}
	aggregate := true
	aggSigScheme := GetAggregateSignatureScheme(clientSignatureScheme, total, batchSize)
	if aggSigScheme == nil {
		aggregate = false
	}
	if aggregate {
		var wg sync.WaitGroup
		for t := 0; t < numBatches; t++ {
			wg.Add(1)
			go func(bn int) {
				start := bn * batchSize
				for i := 0; i < batchSize; i++ {
					aggSigScheme.Aggregate(sigSchemes[start+i], start+i, msgSignatures[start+i], msgHashes[start+i])
				}
				wg.Done()
			}(t)
		}
		wg.Wait()
		result, err := aggSigScheme.Verify()
		if err != nil {
			t.Fatal(err)
		}
		if !result {
			t.Error("signature verification failed")
		}
	} else {
		var wg sync.WaitGroup
		for tr := 0; tr < numBatches; tr++ {
			wg.Add(1)
			go func(bn int) {
				start := bn * batchSize
				for i := 0; i < batchSize; i++ {
					result, err := sigSchemes[start+i].Verify(msgSignatures[start+i], msgHashes[start+i])
					if err != nil {
						t.Fatal(err)
					}
					if !result {
						t.Error("signature verification failed")
					}
				}
				wg.Done()
			}(tr)
		}
		wg.Wait()
	}
}

func TestNewBLS0ChainAggregateSignature(t *testing.T) {
	type args struct {
		total     int
		batchSize int
	}
	tests := []struct {
		name string
		args args
		want *BLS0ChainAggregateSignatureScheme
	}{
		{
			name: "Test_NewBLS0ChainAggregateSignature_OK",
			args: args{total: 1, batchSize: 2},
			want: &BLS0ChainAggregateSignatureScheme{
				Total:     1,
				BatchSize: 2,
				ASigs:     make([]*bls.Sign, 1),
				AGt:       make([]*bls.GT, 1),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewBLS0ChainAggregateSignature(tt.args.total, tt.args.batchSize); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewBLS0ChainAggregateSignature() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBLS0ChainAggregateSignatureScheme_Aggregate(t *testing.T) {
	scheme := NewBLS0ChainScheme()
	if err := scheme.GenerateKeys(); err != nil {
		t.Fatal(err)
	}

	hash := Hash("data")
	sign, err := scheme.Sign(hash)
	if err != nil {
		t.Fatal(err)
	}

	type fields struct {
		Total     int
		BatchSize int
		ASigs     []*bls.Sign
		AGt       []*bls.GT
	}
	type args struct {
		ss        SignatureScheme
		idx       int
		signature string
		hash      string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "Test_BLS0ChainAggregateSignatureScheme_Aggregate_Invalid_Signature_Scheme_ERR",
			args:    args{ss: &ED25519Scheme{}},
			wantErr: true,
		},
		{
			name:    "Test_BLS0ChainAggregateSignatureScheme_Aggregate_Get_Signature_ERR",
			args:    args{ss: &BLS0ChainScheme{}, signature: ""},
			wantErr: true,
		},
		{
			name:    "Test_BLS0ChainAggregateSignatureScheme_Aggregate_Hex_Decoding_Hash_ERR",
			fields:  fields{BatchSize: 1, ASigs: make([]*bls.Sign, 2)},
			args:    args{ss: scheme, signature: sign, hash: "!", idx: 1},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b0a := BLS0ChainAggregateSignatureScheme{
				Total:     tt.fields.Total,
				BatchSize: tt.fields.BatchSize,
				ASigs:     tt.fields.ASigs,
				AGt:       tt.fields.AGt,
			}
			if err := b0a.Aggregate(tt.args.ss, tt.args.idx, tt.args.signature, tt.args.hash); (err != nil) != tt.wantErr {
				t.Errorf("Aggregate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
