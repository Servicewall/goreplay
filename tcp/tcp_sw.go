package tcp

import "fmt"

func (sig *TcpSig) getSYNFp(pckt *Packet) string {
	sig.IPVer = -1
	if pckt.Version == 4 {
		sig.IPVer = 4
	}
	if pckt.Version == 6 {
		sig.IPVer = 6
	}
	return fmt.Sprintf("%[1]d:%s", sig.IPVer, "test")
}

//SW: type from Ganon
type TcpSig struct {
	IPVer   int8
	TTL     uint8
	MSS     int32
	WinType uint8
	WScale  int16
}
