package tcp

import "encoding/binary"

// TCPSynOptions is used to return the parsed TCP Options in a syn
// segment.
type TCPSynOptions struct {
	// MSS is the maximum segment size provided by the peer in the SYN.
	MSS uint16

	// WS is the window scale option provided by the peer in the SYN.
	//
	// Set to -1 if no window scale option was provided.
	WS int

	// TS is true if the timestamp option was provided in the syn/syn-ack.
	TS bool

	// TSVal is the value of the TSVal field in the timestamp option.
	TSVal uint32

	// TSEcr is the value of the TSEcr field in the timestamp option.
	TSEcr uint32

	// SACKPermitted is true if the SACK option was provided in the SYN/SYN-ACK.
	SACKPermitted bool

	OptStr string
}

const (
	TCPOptionEOL           = 0
	TCPOptionNOP           = 1
	TCPOptionMSS           = 2
	TCPOptionWS            = 3
	TCPOptionTS            = 8
	TCPOptionSACKPermitted = 4
	TCPOptionSACK          = 5
)

const (
	// MaxWndScale is maximum allowed window scaling, as described in
	// RFC 1323, section 2.3, page 11.
	MaxWndScale = 14

	// TCPMaxSACKBlocks is the maximum number of SACK blocks that can
	// be encoded in a TCP option field.
	TCPMaxSACKBlocks = 4
)

const (
	// TCPMinimumSize is the minimum size of a valid TCP packet.
	TCPMinimumSize = 20

	// TCPOptionsMaximumSize is the maximum size of TCP options.
	TCPOptionsMaximumSize = 40

	// TCPHeaderMaximumSize is the maximum header size of a TCP packet.
	TCPHeaderMaximumSize = TCPMinimumSize + TCPOptionsMaximumSize

	// TCPProtocolNumber is TCP's transport protocol number.
	//TCPProtocolNumber tcpip.TransportProtocolNumber = 6

	// TCPMinimumMSS is the minimum acceptable value for MSS. This is the
	// same as the value TCP_MIN_MSS defined net/tcp.h.
	//TCPMinimumMSS = IPv4MaximumHeaderSize + TCPHeaderMaximumSize + MinIPFragmentPayloadSize - IPv4MinimumSize - TCPMinimumSize

	// TCPMaximumMSS is the maximum acceptable value for MSS.
	TCPMaximumMSS = 0xffff

	// TCPDefaultMSS is the MSS value that should be used if an MSS option
	// is not received from the peer. It's also the value returned by
	// TCP_MAXSEG option for a socket in an unconnected state.
	//
	// Per RFC 1122, page 85: "If an MSS option is not received at
	// connection setup, TCP MUST assume a default send MSS of 536."
	TCPDefaultMSS = 536
)

func ParseSynOptions(opts []byte, isAck bool) TCPSynOptions {
	limit := len(opts)

	synOpts := TCPSynOptions{
		// Per RFC 1122, page 85: "If an MSS option is not received at
		// connection setup, TCP MUST assume a default send MSS of 536."
		MSS: TCPDefaultMSS,
		// If no window scale option is specified, WS in options is
		// returned as -1; this is because the absence of the option
		// indicates that the we cannot use window scaling on the
		// receive end either.
		WS: -1,
	}

	for i := 0; i < limit; {
		switch opts[i] {
		case TCPOptionEOL:
			i = limit
			//SW
			synOpts.OptStr += "eol+1" //TBD padding
		case TCPOptionNOP:
			i++
			//SW
			synOpts.OptStr += "nop,"

		case TCPOptionMSS:
			if i+4 > limit || opts[i+1] != 4 {
				return synOpts
			}
			mss := uint16(opts[i+2])<<8 | uint16(opts[i+3])
			if mss == 0 {
				return synOpts
			}
			synOpts.MSS = mss
			i += 4
			//SW
			synOpts.OptStr += "mss,"

		case TCPOptionWS:
			if i+3 > limit || opts[i+1] != 3 {
				return synOpts
			}
			ws := int(opts[i+2])
			if ws > MaxWndScale {
				ws = MaxWndScale
			}
			synOpts.WS = ws
			i += 3
			//SW
			synOpts.OptStr += "ws,"

		case TCPOptionTS:
			if i+10 > limit || opts[i+1] != 10 {
				return synOpts
			}
			synOpts.TSVal = binary.BigEndian.Uint32(opts[i+2:])
			if isAck {
				// If the segment is a SYN-ACK then store the Timestamp Echo Reply
				// in the segment.
				synOpts.TSEcr = binary.BigEndian.Uint32(opts[i+6:])
			}
			synOpts.TS = true
			i += 10
			//SW
			synOpts.OptStr += "ts,"

		case TCPOptionSACKPermitted:
			if i+2 > limit || opts[i+1] != 2 {
				return synOpts
			}
			synOpts.SACKPermitted = true
			i += 2
			//SW
			synOpts.OptStr += "sok,"

		default:
			// We don't recognize this option, just skip over it.
			if i+2 > limit {
				return synOpts
			}
			l := int(opts[i+1])
			// If the length is incorrect or if l+i overflows the
			// total options length then return false.
			if l < 2 || i+l > limit {
				return synOpts
			}
			i += l
		}
	}

	return synOpts
}
