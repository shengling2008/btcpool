package eth

import (
	"strings"
	"sync"
	"testing"

	"github.com/Shopify/sarama"
	"github.com/stretchr/testify/assert"
)

type mockProducer struct {
	input chan *sarama.ProducerMessage
}

func (p *mockProducer) AsyncClose() {
}

func (p *mockProducer) Close() error {
	return nil
}

func (p *mockProducer) Input() chan<- *sarama.ProducerMessage {
	return p.input
}

func (p *mockProducer) Successes() <-chan *sarama.ProducerMessage {
	return nil
}

func (p *mockProducer) Errors() <-chan *sarama.ProducerError {
	return nil
}

func TestProcessWork(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)
	work := [10]string{
		"0x070e5f38ffe7a099f7be80433862331e815dde47af4066387f93eebe4d1c0b70",
		"0x45113dd8b32c6ad80229b2576547098555fe464645dca3ad9bd2f9517fc968dd",
		"0x000000000000210a57a1c8cdb7fad922e76fb9a25d0498959dab20b64fe84bb4",
		"0x80a085",
		"0x8cb03169b491f18692502d3b4e2db1c06d6e049b9469491da2077f159046c187",
		"0x7a1200",
		"0x79f3a1",
		"0x92",
		"0x0",
		"0xf901f4a08cb03169b491f18692502d3b4e2db1c06d6e049b9469491da2077f159046c187a01dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d4934794eea5b82b61424df8020f5fedd81767f2d0d25bfba04d2cdbd841680ddd65ca0c6f207312c22a517c115fe87d1124cb6f42094d045ea0bcf838c8b78765abc6fd3c8880a90770dcc61f63feaf6e91c33e3a5ac5d7c6c2a0b41c15c84eb95c3fe1762e42ef19e01f3a0bef93676adf04e8834cfe9be9f979b90100120642c59244144d019c41862184153229829480f5318c29a008c0c082031cb07c0800891888050a80736a86277241cc4a4104042c022403600481c20018c160f200d046c30022805c51300c04080ca468dee82d0980105004b0a90452402220e0d805019318428304bcc187330a8c71a5223240222278a40f04c9904b00615002b09400403080a1049016c1850605d28301f4f18548805440520020f25420b789882004408c6404000d49c088c4acf4630f61446d8480b442a3e0207a2505563090d9278902d01102e40012043d6914044020b18c09198c404c50025096b2020000a68adc9681a01ca20c223188200912358051337aa09b1c181c86a077ce828707bf82d8806eb88380a085837a12008379f3a1845d649dde9dd883010902846765746888676f312e31322e39856c696e757800000000",
	}
	var producer = mockProducer{
		input: make(chan *sarama.ProducerMessage, 1),
	}
	processWork(&wg, work, &producer)
	msg := <-producer.input
	rawGwBytes, _ := msg.Value.Encode()
	rawGw := string(rawGwBytes)
	expected := `"chainType":"ETH","rpcAddress":"","rpcUserPwd":"","parent":"0x8cb03169b491f18692502d3b4e2db1c06d6e049b9469491da2077f159046c187","target":"0x000000000000210a57a1c8cdb7fad922e76fb9a25d0498959dab20b64fe84bb4","hHash":"0x070e5f38ffe7a099f7be80433862331e815dde47af4066387f93eebe4d1c0b70","sHash":"0x45113dd8b32c6ad80229b2576547098555fe464645dca3ad9bd2f9517fc968dd","height":8429701,"uncles":0,"transactions":146,"gasUsedPercent":99.902813,"header":"0xf901f4a08cb03169b491f18692502d3b4e2db1c06d6e049b9469491da2077f159046c187a01dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d4934794eea5b82b61424df8020f5fedd81767f2d0d25bfba04d2cdbd841680ddd65ca0c6f207312c22a517c115fe87d1124cb6f42094d045ea0bcf838c8b78765abc6fd3c8880a90770dcc61f63feaf6e91c33e3a5ac5d7c6c2a0b41c15c84eb95c3fe1762e42ef19e01f3a0bef93676adf04e8834cfe9be9f979b90100120642c59244144d019c41862184153229829480f5318c29a008c0c082031cb07c0800891888050a80736a86277241cc4a4104042c022403600481c20018c160f200d046c30022805c51300c04080ca468dee82d0980105004b0a90452402220e0d805019318428304bcc187330a8c71a5223240222278a40f04c9904b00615002b09400403080a1049016c1850605d28301f4f18548805440520020f25420b789882004408c6404000d49c088c4acf4630f61446d8480b442a3e0207a2505563090d9278902d01102e40012043d6914044020b18c09198c404c50025096b2020000a68adc9681a01ca20c223188200912358051337aa09b1c181c86a077ce828707bf82d8806eb88380a085837a12008379f3a1845d649dde9dd883010902846765746888676f312e31322e39856c696e757800000000"}`
	assert.Equal(t, expected, rawGw[strings.Index(rawGw, `"chainType"`):])
}
