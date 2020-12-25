package tcp

// #include "HelloWorld.h"
import "C"

import (
	"fmt"
	"github.com/buger/goreplay/capture"
	"strings"
)

func (sig *TcpSig) getSYNFp(packet *capture.Packet) string {
	sig.IPVer = -1
	if (packet.NetLayer[0] >> 4) == 4 {
		sig.IPVer = 4
		sig.TTL = packet.NetLayer[8]                    //TTL
		sig.IPOptLen = (packet.NetLayer[0]&0x0F)*4 - 20 //Header Length - 20
		//SYN Package only here
		x := ParseSynOptions(packet.TransLayer[20:20+packet.TransOptsLen], true)
		//println("MSS:", x.MSS, x.WS, x.SACKPermitted, packet.TransOptsLen, packet.TransLayer[20])
		sig.MSS = int32(x.MSS)
		sig.OptStr = x.OptStr
		if strings.Contains(x.OptStr, "mss,nop,ws,nop,nop,ts,sok,eol") {
			sig.Os = "mac"
		} else if strings.Contains(x.OptStr, "mss,nop,ws,nop,nop,sok") {
			sig.Os = "win"
		} else if strings.Contains(x.OptStr, "mss,sok,ts,nop,ws") {
			sig.Os = "lnx"
		}
	} else {
		sig.IPVer = 6
		sig.TTL = packet.NetLayer[7] //TTL
	}

	//C call
	//fmt.Println(C.GoString(C.hello_world()))

	return fmt.Sprintf("%d:%d+%d:%d:%d:%s:%d:%d:%s%s", sig.IPVer, sig.TTL, guessDist(sig.TTL),
		sig.IPOptLen, sig.MSS, sig.OptStr, sig.WinType, sig.WScale, "os->", sig.Os)
}
func guessDist(ttl uint8) uint8 {
	if ttl <= 32 {
		return 32 - ttl
	}
	if ttl <= 64 {
		return 64 - ttl
	}
	if ttl <= 128 {
		return 128 - ttl
	}
	return 255 - ttl
}

//SW: type from Ganon
type TcpSig struct {
	IPVer    int8  /* -1 = any, IP_VER4, IP_VER6         */
	TTL      uint8 /* Actual TTL                         */
	IPOptLen uint8 /* Length of IP options               */
	MSS      int32 /* Maximum segment size (-1 = any)    */
	WinType  uint8 /* WIN_TYPE_*                         */
	WScale   int16 /* Window scale (-1 = any)            */
	OptStr   string
	Os       string /* mac/win/lnx */
}
