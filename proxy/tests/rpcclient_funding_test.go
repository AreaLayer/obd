package main

import (
	"context"
	proxy "github.com/omnilaboratory/obd/proxy/pb"
	"log"
	"testing"
)

func TestOpenChannel(t *testing.T) {

	client, conn := getClient()
	defer conn.Close()

	channelResponse, err := client.OpenChannel(context.Background(), &proxy.OpenChannelRequest{
		RecipientInfo: &proxy.RecipientNodeInfo{
			RecipientNodePeerId: "QmccE4s2uhEXrJXE778NChn1ed8NyWNyAHH23mP7f9NM3L",
			RecipientUserPeerId: "63167817c979ade9e42f3204404c1513a4b1b4e9eea654c9498ed9cc920dbb36"},
		NodePubkeyString: "03c384b8d9c65edea28ce205537bb58dc0096bc618e9e553455e1db1f36cc25642",
		NodePubkeyIndex:  1,
		Private:          false,
	})
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(channelResponse.TemplateChannelId)
}

func TestFundChannel(t *testing.T) {

	client, conn := getClient()
	defer conn.Close()

	fundChannel, err := client.FundChannel(context.Background(), &proxy.FundChannelRequest{
		RecipientInfo: &proxy.RecipientNodeInfo{
			RecipientNodePeerId: "QmccE4s2uhEXrJXE778NChn1ed8NyWNyAHH23mP7f9NM3L",
			RecipientUserPeerId: "63167817c979ade9e42f3204404c1513a4b1b4e9eea654c9498ed9cc920dbb36"},
		TemplateChannelId: "b253ae248c0590e62678f822a25bf9a9e068a2cce92917523c16ed6c7f45d5c5",
		BtcAmount:         0.0004,
		PropertyId:        137,
		AssetAmount:       1,
	})
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(fundChannel.ChannelId)
}