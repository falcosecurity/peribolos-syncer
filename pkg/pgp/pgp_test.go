// Copyright 2023 The Falco Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package pgp_test

import (
	"io"
	"strings"

	"github.com/ProtonMail/go-crypto/openpgp"
	"github.com/ProtonMail/go-crypto/openpgp/armor"
	"github.com/ProtonMail/go-crypto/openpgp/packet"
	"github.com/go-git/go-billy/v5/memfs"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/falcosecurity/peribolos-syncer/pkg/pgp"
)

var _ = Describe("Decoding a PGP private key", func() {
	var (
		err  error
		priv *packet.PrivateKey
	)

	When("the private key is empty", func() {
		BeforeEach(func() {
			priv, err = pgp.DecodePrivateKey(strings.NewReader(""))
		})

		It("should fail", func() {
			Expect(err).ToNot(BeNil())
		})
		It("key be nil", func() {
			Expect(priv).To(BeNil())
		})
	})

	When("the private key is a valid PGP armored private key", func() {
		BeforeEach(func() {
			e, _ := openpgp.NewEntity("", "", "", nil)

			filesystem := memfs.New()
			file, _ := filesystem.Create("mykey")

			armored, _ := armor.Encode(file, openpgp.PrivateKeyType, nil)
			defer armored.Close()

			e.SerializePrivate(armored, nil)
			file.Seek(0, io.SeekStart)

			priv, err = pgp.DecodePrivateKey(file)
		})

		It("should not fail", func() {
			Expect(err).To(BeNil())
		})
		It("should not not be nil", func() {
			Expect(priv).ToNot(BeNil())
		})
	})
})

var _ = Describe("Decoding a PGP public key", func() {
	var (
		err error
		pub *packet.PublicKey
	)

	When("the public key is empty", func() {
		BeforeEach(func() {
			pub, err = pgp.DecodePublicKey(strings.NewReader(""))
		})

		It("should fail", func() {
			Expect(err).ToNot(BeNil())
		})
		It("key be nil", func() {
			Expect(pub).To(BeNil())
		})
	})

	When("the public key is a valid PGP armored public key", func() {
		BeforeEach(func() {
			e, _ := openpgp.NewEntity("", "", "", nil)

			filesystem := memfs.New()
			file, _ := filesystem.Create("mykey")

			armored, _ := armor.Encode(file, openpgp.PublicKeyType, nil)
			defer armored.Close()

			e.Serialize(armored)
			file.Seek(0, io.SeekStart)

			pub, err = pgp.DecodePublicKey(file)
		})

		It("should not fail", func() {
			Expect(err).To(BeNil())
		})
		It("should not not be nil", func() {
			Expect(pub).ToNot(BeNil())
		})
	})
})
