// Written in 2011-2012 by Dmitry Chestnykh.
//
// To the extent possible under law, the author have dedicated all copyright
// and related and neighboring rights to this software to the public domain
// worldwide. This software is distributed without any warranty.
// http://creativecommons.org/publicdomain/zero/1.0/

package blake256

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"hash"
	"testing"
)

func Test256C(t *testing.T) {
	// Test as in C program.
	var hashes = [][]byte{
		{
			0x0C, 0xE8, 0xD4, 0xEF, 0x4D, 0xD7, 0xCD, 0x8D,
			0x62, 0xDF, 0xDE, 0xD9, 0xD4, 0xED, 0xB0, 0xA7,
			0x74, 0xAE, 0x6A, 0x41, 0x92, 0x9A, 0x74, 0xDA,
			0x23, 0x10, 0x9E, 0x8F, 0x11, 0x13, 0x9C, 0x87,
		},
		{
			0xD4, 0x19, 0xBA, 0xD3, 0x2D, 0x50, 0x4F, 0xB7,
			0xD4, 0x4D, 0x46, 0x0C, 0x42, 0xC5, 0x59, 0x3F,
			0xE5, 0x44, 0xFA, 0x4C, 0x13, 0x5D, 0xEC, 0x31,
			0xE2, 0x1B, 0xD9, 0xAB, 0xDC, 0xC2, 0x2D, 0x41,
		},
	}
	data := make([]byte, 72)

	h := New()
	h.Write(data[:1])
	sum := h.Sum(nil)
	if !bytes.Equal(hashes[0], sum) {
		t.Errorf("0: expected %X, got %X", hashes[0], sum)
	}

	// Try to continue hashing.
	h.Write(data[1:])
	sum = h.Sum(nil)
	if !bytes.Equal(hashes[1], sum) {
		t.Errorf("1(1): expected %X, got %X", hashes[1], sum)
	}

	// Try with reset.
	h.Reset()
	h.Write(data)
	sum = h.Sum(nil)
	if !bytes.Equal(hashes[1], sum) {
		t.Errorf("1(2): expected %X, got %X", hashes[1], sum)
	}
}

type blakeVector struct {
	out, in string
}

var vectors256 = []blakeVector{
	{"7576698ee9cad30173080678e5965916adbb11cb5245d386bf1ffda1cb26c9d7",
		"The quick brown fox jumps over the lazy dog"},
	{"07663e00cf96fbc136cf7b1ee099c95346ba3920893d18cc8851f22ee2e36aa6",
		"BLAKE"},
	{"716f6e863f744b9ac22c97ec7b76ea5f5908bc5b2f67c61510bfc4751384ea7a",
		""},
	{"18a393b4e62b1887a2edf79a5c5a5464daf5bbb976f4007bea16a73e4c1e198e",
		"'BLAKE wins SHA-3! Hooray!!!' (I have time machine)"},
	{"fd7282ecc105ef201bb94663fc413db1b7696414682090015f17e309b835f1c2",
		"Go"},
	{"1e75db2a709081f853c2229b65fd1558540aa5e7bd17b04b9a4b31989effa711",
		"HELP! I'm trapped in hash!"},
	{"4181475cb0c22d58ae847e368e91b4669ea2d84bcd55dbf01fe24bae6571dd08",
		`Lorem ipsum dolor sit amet, consectetur adipiscing elit. Donec a diam lectus. Sed sit amet ipsum mauris. Maecenas congue ligula ac quam viverra nec consectetur ante hendrerit. Donec et mollis dolor. Praesent et diam eget libero egestas mattis sit amet vitae augue. Nam tincidunt congue enim, ut porta lorem lacinia consectetur. Donec ut libero sed arcu vehicula ultricies a non tortor. Lorem ipsum dolor sit amet, consectetur adipiscing elit. Aenean ut gravida lorem. Ut turpis felis, pulvinar a semper sed, adipiscing id dolor. Pellentesque auctor nisi id magna consequat sagittis. Curabitur dapibus enim sit amet elit pharetra tincidunt feugiat nisl imperdiet. Ut convallis libero in urna ultrices accumsan. Donec sed odio eros. Donec viverra mi quis quam pulvinar at malesuada arcu rhoncus. Cum sociis natoque penatibus et magnis dis parturient montes, nascetur ridiculus mus. In rutrum accumsan ultricies. Mauris vitae nisi at sem facilisis semper ac in est.`,
	},
	{"af95fffc7768821b1e08866a2f9f66916762bfc9d71c4acb5fd515f31fd6785a", // test with one padding byte
		"Lorem ipsum dolor sit amet, consectetur adipiscing elit. Donec a diam lectus. Sed sit amet ipsum mauris. Maecenas congu",
	},
}

var vectors224 = []blakeVector{
	{"c8e92d7088ef87c1530aee2ad44dc720cc10589cc2ec58f95a15e51b",
		"The quick brown fox jumps over the lazy dog"},
	{"cfb6848add73e1cb47994c4765df33b8f973702705a30a71fe4747a3",
		"BLAKE"},
	{"7dc5313b1c04512a174bd6503b89607aecbee0903d40a8a569c94eed",
		""},
	{"dde9e442003c24495db607b17e07ec1f67396cc1907642a09a96594e",
		"Go"},
	{"9f655b0a92d4155754fa35e055ce7c5e18eb56347081ea1e5158e751",
		"Buffalo buffalo Buffalo buffalo buffalo buffalo Buffalo buffalo"},
}

func testVectors(t *testing.T, hashfunc func() hash.Hash, vectors []blakeVector) {
	for i, v := range vectors {
		h := hashfunc()
		h.Write([]byte(v.in))
		res := fmt.Sprintf("%x", h.Sum(nil))
		if res != v.out {
			t.Errorf("%d: expected %q, got %q", i, v.out, res)
		}
	}
}

func Test256(t *testing.T) {
	testVectors(t, New, vectors256)
}

func Test224(t *testing.T) {
	testVectors(t, New224, vectors224)
}

var vectors256salt = []struct{ out, in, salt string }{
	{"561d6d0cfa3d31d5eedaf2d575f3942539b03522befc2a1196ba0e51af8992a8",
		"",
		"1234567890123456"},
	{"88cc11889bbbee42095337fe2153c591971f94fbf8fe540d3c7e9f1700ab2d0c",
		"It's so salty out there!",
		"SALTsaltSaltSALT"},
}

func TestSalt(t *testing.T) {
	for i, v := range vectors256salt {
		h := NewSalt([]byte(v.salt))
		h.Write([]byte(v.in))
		res := fmt.Sprintf("%x", h.Sum(nil))
		if res != v.out {
			t.Errorf("%d: expected %q, got %q", i, v.out, res)
		}
	}

	// Check that passing bad salt length panics.
	defer func() {
		if err := recover(); err == nil {
			t.Errorf("expected panic for bad salt length")
		}
	}()
	NewSalt([]byte{1, 2, 3, 4, 5, 6, 7, 8})
}

var longerData, longData, shortData []byte

func init() {
	longerData = make([]byte, 88888)
	longData = make([]byte, 4096)
	shortData = make([]byte, 64)
}

func benchHash(b *testing.B, hashfunc func() hash.Hash, data []byte) {
	b.StopTimer()
	h := hashfunc()
	digest := make([]byte, 0, BlockSize)
	b.SetBytes(int64(len(data)))
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		h.Write(data)
		h.Sum(digest)
		h.Reset()
	}
}

func benchHashWrite(b *testing.B, hashfunc func() hash.Hash, data []byte) {
	b.StopTimer()
	h := hashfunc()
	b.SetBytes(int64(len(data)))
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		h.Write(data)
	}
}

func BenchmarkLonger(b *testing.B) {
	benchHash(b, New, longerData)
}

func BenchmarkLong(b *testing.B) {
	benchHash(b, New, longData)
}

func BenchmarkShort(b *testing.B) {
	benchHash(b, New, shortData)
}

func BenchmarkSHA2LL(b *testing.B) {
	benchHash(b, sha256.New, longerData)
}

func BenchmarkSHA2L(b *testing.B) {
	benchHash(b, sha256.New, longData)
}

func BenchmarkSHA2S(b *testing.B) {
	benchHash(b, sha256.New, shortData)
}

func BenchmarkWriteLonger(b *testing.B) {
	benchHashWrite(b, New, longerData)
}

func BenchmarkWriteLong(b *testing.B) {
	benchHashWrite(b, New, longData)
}

func BenchmarkWriteShort(b *testing.B) {
	benchHashWrite(b, New, shortData)
}

func BenchmarkWriteSHA2LL(b *testing.B) {
	benchHashWrite(b, sha256.New, longerData)
}

func BenchmarkWriteSHA2L(b *testing.B) {
	benchHashWrite(b, sha256.New, longData)
}

func BenchmarkWriteSHA2S(b *testing.B) {
	benchHashWrite(b, sha256.New, shortData)
}
