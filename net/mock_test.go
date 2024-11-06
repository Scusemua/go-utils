package net_test

import (
	"fmt"
	"net"

	mock "github.com/jordwest/mock-conn"

	. "github.com/Scusemua/go-utils/net"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var (
	shortcutAddress = "shortcut:%d:%s"
)

var _ = Describe("Mock", func() {
	newResponseConn := func() (net.Conn, net.Conn) {
		conn := mock.NewConn()
		return conn.Client, conn.Server
	}

	It("should mock connection works", func() {
		buff := make([]byte, 1024)
		client, server := newResponseConn()
		// reqReader, reqWriter, _, _ := newResponseConn()
		go func() {
			n, err := client.Write([]byte("ping"))
			Expect(n).To(Equal(4))
			Expect(err).To(BeNil())
		}()

		read, err := server.Read(buff)
		Expect(err).To(BeNil())
		Expect(string(buff[:read])).To(Equal("ping"))

		go func() {
			n, err := server.Write([]byte("pong"))
			Expect(n).To(Equal(4))
			Expect(err).To(BeNil())
		}()

		read, err = client.Read(buff)
		Expect(err).To(BeNil())
		Expect(string(buff[:read])).To(Equal("pong"))
	})

	It("should Validate recognize shortcut", func() {
		shortcut := InitShortcut()
		ip := "10.23.4.5"
		fullip := ip + ":6378"

		addr, ok := shortcut.Validate(fmt.Sprintf(shortcutAddress, 1, ip))
		Expect(ok).To(Equal(true))
		Expect(addr).To(Equal(ip))

		addr, ok = shortcut.Validate(fmt.Sprintf(shortcutAddress, 1, fullip))
		Expect(ok).To(Equal(true))
		Expect(addr).To(Equal(fullip))

		addr, ok = shortcut.Validate(fullip)
		Expect(ok).To(Equal(false))
		Expect(addr).To(Equal(""))
	})
})
