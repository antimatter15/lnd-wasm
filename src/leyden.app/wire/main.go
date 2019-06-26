package main

import (
	"fmt"
	// "syscall/js"
	// "encoding/hex"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/btcec"
	// "time"
	"bytes"
	"net"
	"encoding/hex"
	// "github.com/antimatter15/lnwasm/net"
	"github.com/lightningnetwork/lnd/lnwire"
	"github.com/lightningnetwork/lnd/brontide"
	// "sync"
	// "github.com/antimatter15/lnwasm/lnwire"
	"strconv"
	"image/color"


)

type WumboAddr struct {
	Host string
	Port int
}



// A compile-time check to ensure that OnionAddr implements the net.Addr
// interface.
var _ net.Addr = (*WumboAddr)(nil)

// String returns the string representation of an onion address.
func (o *WumboAddr) String() string {
	return net.JoinHostPort(o.Host, strconv.Itoa(o.Port))
}

// Network returns the network that this implementation of net.Addr will use.
// In this case, because Tor only allows TCP connections, the network is "tcp".
func (o *WumboAddr) Network() string {
	return "tcp"
}


func main() {
	fmt.Println("Hello, WebAssembly stuff!")

	localPriv, err := btcec.NewPrivateKey(btcec.S256())
	if err != nil {
		fmt.Println("unable to make pkey")
		return
	}

	// netAddr := lnwire.NetAddress{
	// 	IdentityKey
	// }

	// "03c2abfa93eacec04721c019644584424aab2ba4dff3ac9bdab4e9c97007491dda@104.248.84.249:9735"
	parsedPubKey := "03c2abfa93eacec04721c019644584424aab2ba4dff3ac9bdab4e9c97007491dda"

	pubKeyBytes, err := hex.DecodeString(parsedPubKey)
	if err != nil {
		fmt.Println("invalid lightning address pubkey: %v", err)
		return 
	}

	// The compressed pubkey should have a length of exactly 33 bytes.
	if len(pubKeyBytes) != 33 {
		fmt.Println("invalid lightning address pubkey: "+
			"length must be 33 bytes, found %d", len(pubKeyBytes))
		return 
	}

	// Parse the pubkey bytes to verify that it corresponds to valid public
	// key on the secp256k1 curve.
	pubKey, err := btcec.ParsePubKey(pubKeyBytes, btcec.S256())
	if err != nil {
		fmt.Println("invalid lightning address pubkey: %v", err)
		return 
	}

	addr := &WumboAddr{
		Host: "104.248.84.249",
		Port: 9735,
	}


	fmt.Println(pubKey)
	// return

	netAddr := &lnwire.NetAddress{
		IdentityKey: pubKey,
		Address:     addr,
	}

	fmt.Println(netAddr)

	// aliceTCPAddr, _ = net.ResolveTCPAddr("tcp", "10.0.0.2:9001")

	// aliceAddr = &lnwire.NetAddress{
	// 	IdentityKey: alicePubKey,
	// 	Address:     aliceTCPAddr,
	// }


	thingy := ClearNet{}

	conn, err := brontide.Dial(localPriv, netAddr, thingy.Dial)
	if err != nil {
		fmt.Println("Unable to connect to %v: %v", netAddr, err)
		return
	}

	fmt.Println("connected woot!", conn)


	rawMsg, _ := conn.ReadNextMessage()
	msgReader := bytes.NewReader(rawMsg)
	nextMsg, err := lnwire.ReadMessage(msgReader, 0)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(nextMsg.MsgType().String())


	// fmt.Println(conn.ReadNextMessage())

	alias, _ := lnwire.NewNodeAlias("hello kevin was here awake")

	req := lnwire.NodeAnnouncement{
		Features:  lnwire.NewRawFeatureVector(),
		Timestamp: 42,
		Alias:     alias,
		RGBColor: color.RGBA{
			R: 0,
			G: 127,
			B: 255,
		},
	}
	data, _ := req.DataToSign()
	msgDigest := chainhash.DoubleHashB(data)
	sign, _ := localPriv.Sign(msgDigest)
	req.Signature, err = lnwire.NewSigFromSignature(sign)

	// req.Signature, err = NewSigFromSignature(testSig)
	if err != nil {
		fmt.Println("unable to parse sig: %v", err)
		return
	}

	// req.NodeID, err = randRawKey()
	// if err != nil {
	// 	t.Fatalf("unable to generate key: %v", err)
	// 	return
	// }

	// req.Addresses, err = randAddrs(r)
	// if err != nil {
	// 	t.Fatalf("unable to generate addresses: %v", err)
	// }

	// numExtraBytes := r.Int31n(1000)
	// if numExtraBytes > 0 {
	// 	req.ExtraOpaqueData = make([]byte, numExtraBytes)
	// 	_, err := r.Read(req.ExtraOpaqueData[:])
	// 	if err != nil {
	// 		t.Fatalf("unable to generate opaque "+
	// 			"bytes: %v", err)
	// 		return
	// 	}
	// }

	// v[0] = reflect.ValueOf(req)


	var b bytes.Buffer
	if _, err := lnwire.WriteMessage(&b, &req, 0); err != nil {
		fmt.Println("unable to write msg: %v", err)
		return 
	}

	conn.WriteMessage(b.Bytes())

	// Dial(localPriv *btcec.PrivateKey, netAddr *lnwire.NetAddress,
	// dialer func(string, string) (net.Conn, error))

}
