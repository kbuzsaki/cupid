// ************************************************************
// DO NOT EDIT.
// THIS FILE IS AUTO-GENERATED BY codecgen.
// ************************************************************

package client

import (
	"errors"
	"fmt"
	"reflect"
	"runtime"
	time "time"

	codec1978 "github.com/ugorji/go/codec"
)

const (
	// ----- content types ----
	codecSelferC_UTF81819 = 1
	codecSelferC_RAW1819  = 0
	// ----- value types used ----
	codecSelferValueTypeArray1819 = 10
	codecSelferValueTypeMap1819   = 9
	// ----- containerStateValues ----
	codecSelfer_containerMapKey1819    = 2
	codecSelfer_containerMapValue1819  = 3
	codecSelfer_containerMapEnd1819    = 4
	codecSelfer_containerArrayElem1819 = 6
	codecSelfer_containerArrayEnd1819  = 7
)

var (
	codecSelferBitsize1819                         = uint8(reflect.TypeOf(uint(0)).Bits())
	codecSelferOnlyMapOrArrayEncodeToStructErr1819 = errors.New(`only encoded map or array can be decoded into a struct`)
)

type codecSelfer1819 struct{}

func init() {
	if codec1978.GenVersion != 5 {
		_, file, _, _ := runtime.Caller(0)
		err := fmt.Errorf("codecgen version mismatch: current: %v, need %v. Re-generate file: %v",
			5, codec1978.GenVersion, file)
		panic(err)
	}
	if false { // reference the types, but skip this branch at build/run time
		var v0 time.Time
		_ = v0
	}
}

func (x *Response) CodecEncodeSelf(e *codec1978.Encoder) {
	var h codecSelfer1819
	z, r := codec1978.GenHelperEncoder(e)
	_, _, _ = h, z, r
	if x == nil {
		r.EncodeNil()
	} else {
		yym1 := z.EncBinary()
		_ = yym1
		if false {
		} else if z.HasExtensions() && z.EncExt(x) {
		} else {
			yysep2 := !z.EncBinary()
			yy2arr2 := z.EncBasicHandle().StructToArray
			var yyq2 [3]bool
			_, _, _ = yysep2, yyq2, yy2arr2
			const yyr2 bool = false
			var yynn2 int
			if yyr2 || yy2arr2 {
				r.EncodeArrayStart(3)
			} else {
				yynn2 = 3
				for _, b := range yyq2 {
					if b {
						yynn2++
					}
				}
				r.EncodeMapStart(yynn2)
				yynn2 = 0
			}
			if yyr2 || yy2arr2 {
				z.EncSendContainerState(codecSelfer_containerArrayElem1819)
				yym4 := z.EncBinary()
				_ = yym4
				if false {
				} else {
					r.EncodeString(codecSelferC_UTF81819, string(x.Action))
				}
			} else {
				z.EncSendContainerState(codecSelfer_containerMapKey1819)
				r.EncodeString(codecSelferC_UTF81819, string("action"))
				z.EncSendContainerState(codecSelfer_containerMapValue1819)
				yym5 := z.EncBinary()
				_ = yym5
				if false {
				} else {
					r.EncodeString(codecSelferC_UTF81819, string(x.Action))
				}
			}
			if yyr2 || yy2arr2 {
				z.EncSendContainerState(codecSelfer_containerArrayElem1819)
				if x.Node == nil {
					r.EncodeNil()
				} else {
					x.Node.CodecEncodeSelf(e)
				}
			} else {
				z.EncSendContainerState(codecSelfer_containerMapKey1819)
				r.EncodeString(codecSelferC_UTF81819, string("node"))
				z.EncSendContainerState(codecSelfer_containerMapValue1819)
				if x.Node == nil {
					r.EncodeNil()
				} else {
					x.Node.CodecEncodeSelf(e)
				}
			}
			if yyr2 || yy2arr2 {
				z.EncSendContainerState(codecSelfer_containerArrayElem1819)
				if x.PrevNode == nil {
					r.EncodeNil()
				} else {
					x.PrevNode.CodecEncodeSelf(e)
				}
			} else {
				z.EncSendContainerState(codecSelfer_containerMapKey1819)
				r.EncodeString(codecSelferC_UTF81819, string("prevNode"))
				z.EncSendContainerState(codecSelfer_containerMapValue1819)
				if x.PrevNode == nil {
					r.EncodeNil()
				} else {
					x.PrevNode.CodecEncodeSelf(e)
				}
			}
			if yyr2 || yy2arr2 {
				z.EncSendContainerState(codecSelfer_containerArrayEnd1819)
			} else {
				z.EncSendContainerState(codecSelfer_containerMapEnd1819)
			}
		}
	}
}

func (x *Response) CodecDecodeSelf(d *codec1978.Decoder) {
	var h codecSelfer1819
	z, r := codec1978.GenHelperDecoder(d)
	_, _, _ = h, z, r
	yym1 := z.DecBinary()
	_ = yym1
	if false {
	} else if z.HasExtensions() && z.DecExt(x) {
	} else {
		yyct2 := r.ContainerType()
		if yyct2 == codecSelferValueTypeMap1819 {
			yyl2 := r.ReadMapStart()
			if yyl2 == 0 {
				z.DecSendContainerState(codecSelfer_containerMapEnd1819)
			} else {
				x.codecDecodeSelfFromMap(yyl2, d)
			}
		} else if yyct2 == codecSelferValueTypeArray1819 {
			yyl2 := r.ReadArrayStart()
			if yyl2 == 0 {
				z.DecSendContainerState(codecSelfer_containerArrayEnd1819)
			} else {
				x.codecDecodeSelfFromArray(yyl2, d)
			}
		} else {
			panic(codecSelferOnlyMapOrArrayEncodeToStructErr1819)
		}
	}
}

func (x *Response) codecDecodeSelfFromMap(l int, d *codec1978.Decoder) {
	var h codecSelfer1819
	z, r := codec1978.GenHelperDecoder(d)
	_, _, _ = h, z, r
	var yys3Slc = z.DecScratchBuffer() // default slice to decode into
	_ = yys3Slc
	var yyhl3 bool = l >= 0
	for yyj3 := 0; ; yyj3++ {
		if yyhl3 {
			if yyj3 >= l {
				break
			}
		} else {
			if r.CheckBreak() {
				break
			}
		}
		z.DecSendContainerState(codecSelfer_containerMapKey1819)
		yys3Slc = r.DecodeBytes(yys3Slc, true, true)
		yys3 := string(yys3Slc)
		z.DecSendContainerState(codecSelfer_containerMapValue1819)
		switch yys3 {
		case "action":
			if r.TryDecodeAsNil() {
				x.Action = ""
			} else {
				yyv4 := &x.Action
				yym5 := z.DecBinary()
				_ = yym5
				if false {
				} else {
					*((*string)(yyv4)) = r.DecodeString()
				}
			}
		case "node":
			if r.TryDecodeAsNil() {
				if x.Node != nil {
					x.Node = nil
				}
			} else {
				if x.Node == nil {
					x.Node = new(Node)
				}
				x.Node.CodecDecodeSelf(d)
			}
		case "prevNode":
			if r.TryDecodeAsNil() {
				if x.PrevNode != nil {
					x.PrevNode = nil
				}
			} else {
				if x.PrevNode == nil {
					x.PrevNode = new(Node)
				}
				x.PrevNode.CodecDecodeSelf(d)
			}
		default:
			z.DecStructFieldNotFound(-1, yys3)
		} // end switch yys3
	} // end for yyj3
	z.DecSendContainerState(codecSelfer_containerMapEnd1819)
}

func (x *Response) codecDecodeSelfFromArray(l int, d *codec1978.Decoder) {
	var h codecSelfer1819
	z, r := codec1978.GenHelperDecoder(d)
	_, _, _ = h, z, r
	var yyj8 int
	var yyb8 bool
	var yyhl8 bool = l >= 0
	yyj8++
	if yyhl8 {
		yyb8 = yyj8 > l
	} else {
		yyb8 = r.CheckBreak()
	}
	if yyb8 {
		z.DecSendContainerState(codecSelfer_containerArrayEnd1819)
		return
	}
	z.DecSendContainerState(codecSelfer_containerArrayElem1819)
	if r.TryDecodeAsNil() {
		x.Action = ""
	} else {
		yyv9 := &x.Action
		yym10 := z.DecBinary()
		_ = yym10
		if false {
		} else {
			*((*string)(yyv9)) = r.DecodeString()
		}
	}
	yyj8++
	if yyhl8 {
		yyb8 = yyj8 > l
	} else {
		yyb8 = r.CheckBreak()
	}
	if yyb8 {
		z.DecSendContainerState(codecSelfer_containerArrayEnd1819)
		return
	}
	z.DecSendContainerState(codecSelfer_containerArrayElem1819)
	if r.TryDecodeAsNil() {
		if x.Node != nil {
			x.Node = nil
		}
	} else {
		if x.Node == nil {
			x.Node = new(Node)
		}
		x.Node.CodecDecodeSelf(d)
	}
	yyj8++
	if yyhl8 {
		yyb8 = yyj8 > l
	} else {
		yyb8 = r.CheckBreak()
	}
	if yyb8 {
		z.DecSendContainerState(codecSelfer_containerArrayEnd1819)
		return
	}
	z.DecSendContainerState(codecSelfer_containerArrayElem1819)
	if r.TryDecodeAsNil() {
		if x.PrevNode != nil {
			x.PrevNode = nil
		}
	} else {
		if x.PrevNode == nil {
			x.PrevNode = new(Node)
		}
		x.PrevNode.CodecDecodeSelf(d)
	}
	for {
		yyj8++
		if yyhl8 {
			yyb8 = yyj8 > l
		} else {
			yyb8 = r.CheckBreak()
		}
		if yyb8 {
			break
		}
		z.DecSendContainerState(codecSelfer_containerArrayElem1819)
		z.DecStructFieldNotFound(yyj8-1, "")
	}
	z.DecSendContainerState(codecSelfer_containerArrayEnd1819)
}

func (x *Node) CodecEncodeSelf(e *codec1978.Encoder) {
	var h codecSelfer1819
	z, r := codec1978.GenHelperEncoder(e)
	_, _, _ = h, z, r
	if x == nil {
		r.EncodeNil()
	} else {
		yym1 := z.EncBinary()
		_ = yym1
		if false {
		} else if z.HasExtensions() && z.EncExt(x) {
		} else {
			yysep2 := !z.EncBinary()
			yy2arr2 := z.EncBasicHandle().StructToArray
			var yyq2 [8]bool
			_, _, _ = yysep2, yyq2, yy2arr2
			const yyr2 bool = false
			yyq2[1] = x.Dir != false
			yyq2[6] = x.Expiration != nil
			yyq2[7] = x.TTL != 0
			var yynn2 int
			if yyr2 || yy2arr2 {
				r.EncodeArrayStart(8)
			} else {
				yynn2 = 5
				for _, b := range yyq2 {
					if b {
						yynn2++
					}
				}
				r.EncodeMapStart(yynn2)
				yynn2 = 0
			}
			if yyr2 || yy2arr2 {
				z.EncSendContainerState(codecSelfer_containerArrayElem1819)
				yym4 := z.EncBinary()
				_ = yym4
				if false {
				} else {
					r.EncodeString(codecSelferC_UTF81819, string(x.Key))
				}
			} else {
				z.EncSendContainerState(codecSelfer_containerMapKey1819)
				r.EncodeString(codecSelferC_UTF81819, string("key"))
				z.EncSendContainerState(codecSelfer_containerMapValue1819)
				yym5 := z.EncBinary()
				_ = yym5
				if false {
				} else {
					r.EncodeString(codecSelferC_UTF81819, string(x.Key))
				}
			}
			if yyr2 || yy2arr2 {
				z.EncSendContainerState(codecSelfer_containerArrayElem1819)
				if yyq2[1] {
					yym7 := z.EncBinary()
					_ = yym7
					if false {
					} else {
						r.EncodeBool(bool(x.Dir))
					}
				} else {
					r.EncodeBool(false)
				}
			} else {
				if yyq2[1] {
					z.EncSendContainerState(codecSelfer_containerMapKey1819)
					r.EncodeString(codecSelferC_UTF81819, string("dir"))
					z.EncSendContainerState(codecSelfer_containerMapValue1819)
					yym8 := z.EncBinary()
					_ = yym8
					if false {
					} else {
						r.EncodeBool(bool(x.Dir))
					}
				}
			}
			if yyr2 || yy2arr2 {
				z.EncSendContainerState(codecSelfer_containerArrayElem1819)
				yym10 := z.EncBinary()
				_ = yym10
				if false {
				} else {
					r.EncodeString(codecSelferC_UTF81819, string(x.Value))
				}
			} else {
				z.EncSendContainerState(codecSelfer_containerMapKey1819)
				r.EncodeString(codecSelferC_UTF81819, string("value"))
				z.EncSendContainerState(codecSelfer_containerMapValue1819)
				yym11 := z.EncBinary()
				_ = yym11
				if false {
				} else {
					r.EncodeString(codecSelferC_UTF81819, string(x.Value))
				}
			}
			if yyr2 || yy2arr2 {
				z.EncSendContainerState(codecSelfer_containerArrayElem1819)
				if x.Nodes == nil {
					r.EncodeNil()
				} else {
					x.Nodes.CodecEncodeSelf(e)
				}
			} else {
				z.EncSendContainerState(codecSelfer_containerMapKey1819)
				r.EncodeString(codecSelferC_UTF81819, string("nodes"))
				z.EncSendContainerState(codecSelfer_containerMapValue1819)
				if x.Nodes == nil {
					r.EncodeNil()
				} else {
					x.Nodes.CodecEncodeSelf(e)
				}
			}
			if yyr2 || yy2arr2 {
				z.EncSendContainerState(codecSelfer_containerArrayElem1819)
				yym16 := z.EncBinary()
				_ = yym16
				if false {
				} else {
					r.EncodeUint(uint64(x.CreatedIndex))
				}
			} else {
				z.EncSendContainerState(codecSelfer_containerMapKey1819)
				r.EncodeString(codecSelferC_UTF81819, string("createdIndex"))
				z.EncSendContainerState(codecSelfer_containerMapValue1819)
				yym17 := z.EncBinary()
				_ = yym17
				if false {
				} else {
					r.EncodeUint(uint64(x.CreatedIndex))
				}
			}
			if yyr2 || yy2arr2 {
				z.EncSendContainerState(codecSelfer_containerArrayElem1819)
				yym19 := z.EncBinary()
				_ = yym19
				if false {
				} else {
					r.EncodeUint(uint64(x.ModifiedIndex))
				}
			} else {
				z.EncSendContainerState(codecSelfer_containerMapKey1819)
				r.EncodeString(codecSelferC_UTF81819, string("modifiedIndex"))
				z.EncSendContainerState(codecSelfer_containerMapValue1819)
				yym20 := z.EncBinary()
				_ = yym20
				if false {
				} else {
					r.EncodeUint(uint64(x.ModifiedIndex))
				}
			}
			if yyr2 || yy2arr2 {
				z.EncSendContainerState(codecSelfer_containerArrayElem1819)
				if yyq2[6] {
					if x.Expiration == nil {
						r.EncodeNil()
					} else {
						yym22 := z.EncBinary()
						_ = yym22
						if false {
						} else if yym23 := z.TimeRtidIfBinc(); yym23 != 0 {
							r.EncodeBuiltin(yym23, x.Expiration)
						} else if z.HasExtensions() && z.EncExt(x.Expiration) {
						} else if yym22 {
							z.EncBinaryMarshal(x.Expiration)
						} else if !yym22 && z.IsJSONHandle() {
							z.EncJSONMarshal(x.Expiration)
						} else {
							z.EncFallback(x.Expiration)
						}
					}
				} else {
					r.EncodeNil()
				}
			} else {
				if yyq2[6] {
					z.EncSendContainerState(codecSelfer_containerMapKey1819)
					r.EncodeString(codecSelferC_UTF81819, string("expiration"))
					z.EncSendContainerState(codecSelfer_containerMapValue1819)
					if x.Expiration == nil {
						r.EncodeNil()
					} else {
						yym24 := z.EncBinary()
						_ = yym24
						if false {
						} else if yym25 := z.TimeRtidIfBinc(); yym25 != 0 {
							r.EncodeBuiltin(yym25, x.Expiration)
						} else if z.HasExtensions() && z.EncExt(x.Expiration) {
						} else if yym24 {
							z.EncBinaryMarshal(x.Expiration)
						} else if !yym24 && z.IsJSONHandle() {
							z.EncJSONMarshal(x.Expiration)
						} else {
							z.EncFallback(x.Expiration)
						}
					}
				}
			}
			if yyr2 || yy2arr2 {
				z.EncSendContainerState(codecSelfer_containerArrayElem1819)
				if yyq2[7] {
					yym27 := z.EncBinary()
					_ = yym27
					if false {
					} else {
						r.EncodeInt(int64(x.TTL))
					}
				} else {
					r.EncodeInt(0)
				}
			} else {
				if yyq2[7] {
					z.EncSendContainerState(codecSelfer_containerMapKey1819)
					r.EncodeString(codecSelferC_UTF81819, string("ttl"))
					z.EncSendContainerState(codecSelfer_containerMapValue1819)
					yym28 := z.EncBinary()
					_ = yym28
					if false {
					} else {
						r.EncodeInt(int64(x.TTL))
					}
				}
			}
			if yyr2 || yy2arr2 {
				z.EncSendContainerState(codecSelfer_containerArrayEnd1819)
			} else {
				z.EncSendContainerState(codecSelfer_containerMapEnd1819)
			}
		}
	}
}

func (x *Node) CodecDecodeSelf(d *codec1978.Decoder) {
	var h codecSelfer1819
	z, r := codec1978.GenHelperDecoder(d)
	_, _, _ = h, z, r
	yym1 := z.DecBinary()
	_ = yym1
	if false {
	} else if z.HasExtensions() && z.DecExt(x) {
	} else {
		yyct2 := r.ContainerType()
		if yyct2 == codecSelferValueTypeMap1819 {
			yyl2 := r.ReadMapStart()
			if yyl2 == 0 {
				z.DecSendContainerState(codecSelfer_containerMapEnd1819)
			} else {
				x.codecDecodeSelfFromMap(yyl2, d)
			}
		} else if yyct2 == codecSelferValueTypeArray1819 {
			yyl2 := r.ReadArrayStart()
			if yyl2 == 0 {
				z.DecSendContainerState(codecSelfer_containerArrayEnd1819)
			} else {
				x.codecDecodeSelfFromArray(yyl2, d)
			}
		} else {
			panic(codecSelferOnlyMapOrArrayEncodeToStructErr1819)
		}
	}
}

func (x *Node) codecDecodeSelfFromMap(l int, d *codec1978.Decoder) {
	var h codecSelfer1819
	z, r := codec1978.GenHelperDecoder(d)
	_, _, _ = h, z, r
	var yys3Slc = z.DecScratchBuffer() // default slice to decode into
	_ = yys3Slc
	var yyhl3 bool = l >= 0
	for yyj3 := 0; ; yyj3++ {
		if yyhl3 {
			if yyj3 >= l {
				break
			}
		} else {
			if r.CheckBreak() {
				break
			}
		}
		z.DecSendContainerState(codecSelfer_containerMapKey1819)
		yys3Slc = r.DecodeBytes(yys3Slc, true, true)
		yys3 := string(yys3Slc)
		z.DecSendContainerState(codecSelfer_containerMapValue1819)
		switch yys3 {
		case "key":
			if r.TryDecodeAsNil() {
				x.Key = ""
			} else {
				yyv4 := &x.Key
				yym5 := z.DecBinary()
				_ = yym5
				if false {
				} else {
					*((*string)(yyv4)) = r.DecodeString()
				}
			}
		case "dir":
			if r.TryDecodeAsNil() {
				x.Dir = false
			} else {
				yyv6 := &x.Dir
				yym7 := z.DecBinary()
				_ = yym7
				if false {
				} else {
					*((*bool)(yyv6)) = r.DecodeBool()
				}
			}
		case "value":
			if r.TryDecodeAsNil() {
				x.Value = ""
			} else {
				yyv8 := &x.Value
				yym9 := z.DecBinary()
				_ = yym9
				if false {
				} else {
					*((*string)(yyv8)) = r.DecodeString()
				}
			}
		case "nodes":
			if r.TryDecodeAsNil() {
				x.Nodes = nil
			} else {
				yyv10 := &x.Nodes
				yyv10.CodecDecodeSelf(d)
			}
		case "createdIndex":
			if r.TryDecodeAsNil() {
				x.CreatedIndex = 0
			} else {
				yyv11 := &x.CreatedIndex
				yym12 := z.DecBinary()
				_ = yym12
				if false {
				} else {
					*((*uint64)(yyv11)) = uint64(r.DecodeUint(64))
				}
			}
		case "modifiedIndex":
			if r.TryDecodeAsNil() {
				x.ModifiedIndex = 0
			} else {
				yyv13 := &x.ModifiedIndex
				yym14 := z.DecBinary()
				_ = yym14
				if false {
				} else {
					*((*uint64)(yyv13)) = uint64(r.DecodeUint(64))
				}
			}
		case "expiration":
			if r.TryDecodeAsNil() {
				if x.Expiration != nil {
					x.Expiration = nil
				}
			} else {
				if x.Expiration == nil {
					x.Expiration = new(time.Time)
				}
				yym16 := z.DecBinary()
				_ = yym16
				if false {
				} else if yym17 := z.TimeRtidIfBinc(); yym17 != 0 {
					r.DecodeBuiltin(yym17, x.Expiration)
				} else if z.HasExtensions() && z.DecExt(x.Expiration) {
				} else if yym16 {
					z.DecBinaryUnmarshal(x.Expiration)
				} else if !yym16 && z.IsJSONHandle() {
					z.DecJSONUnmarshal(x.Expiration)
				} else {
					z.DecFallback(x.Expiration, false)
				}
			}
		case "ttl":
			if r.TryDecodeAsNil() {
				x.TTL = 0
			} else {
				yyv18 := &x.TTL
				yym19 := z.DecBinary()
				_ = yym19
				if false {
				} else {
					*((*int64)(yyv18)) = int64(r.DecodeInt(64))
				}
			}
		default:
			z.DecStructFieldNotFound(-1, yys3)
		} // end switch yys3
	} // end for yyj3
	z.DecSendContainerState(codecSelfer_containerMapEnd1819)
}

func (x *Node) codecDecodeSelfFromArray(l int, d *codec1978.Decoder) {
	var h codecSelfer1819
	z, r := codec1978.GenHelperDecoder(d)
	_, _, _ = h, z, r
	var yyj20 int
	var yyb20 bool
	var yyhl20 bool = l >= 0
	yyj20++
	if yyhl20 {
		yyb20 = yyj20 > l
	} else {
		yyb20 = r.CheckBreak()
	}
	if yyb20 {
		z.DecSendContainerState(codecSelfer_containerArrayEnd1819)
		return
	}
	z.DecSendContainerState(codecSelfer_containerArrayElem1819)
	if r.TryDecodeAsNil() {
		x.Key = ""
	} else {
		yyv21 := &x.Key
		yym22 := z.DecBinary()
		_ = yym22
		if false {
		} else {
			*((*string)(yyv21)) = r.DecodeString()
		}
	}
	yyj20++
	if yyhl20 {
		yyb20 = yyj20 > l
	} else {
		yyb20 = r.CheckBreak()
	}
	if yyb20 {
		z.DecSendContainerState(codecSelfer_containerArrayEnd1819)
		return
	}
	z.DecSendContainerState(codecSelfer_containerArrayElem1819)
	if r.TryDecodeAsNil() {
		x.Dir = false
	} else {
		yyv23 := &x.Dir
		yym24 := z.DecBinary()
		_ = yym24
		if false {
		} else {
			*((*bool)(yyv23)) = r.DecodeBool()
		}
	}
	yyj20++
	if yyhl20 {
		yyb20 = yyj20 > l
	} else {
		yyb20 = r.CheckBreak()
	}
	if yyb20 {
		z.DecSendContainerState(codecSelfer_containerArrayEnd1819)
		return
	}
	z.DecSendContainerState(codecSelfer_containerArrayElem1819)
	if r.TryDecodeAsNil() {
		x.Value = ""
	} else {
		yyv25 := &x.Value
		yym26 := z.DecBinary()
		_ = yym26
		if false {
		} else {
			*((*string)(yyv25)) = r.DecodeString()
		}
	}
	yyj20++
	if yyhl20 {
		yyb20 = yyj20 > l
	} else {
		yyb20 = r.CheckBreak()
	}
	if yyb20 {
		z.DecSendContainerState(codecSelfer_containerArrayEnd1819)
		return
	}
	z.DecSendContainerState(codecSelfer_containerArrayElem1819)
	if r.TryDecodeAsNil() {
		x.Nodes = nil
	} else {
		yyv27 := &x.Nodes
		yyv27.CodecDecodeSelf(d)
	}
	yyj20++
	if yyhl20 {
		yyb20 = yyj20 > l
	} else {
		yyb20 = r.CheckBreak()
	}
	if yyb20 {
		z.DecSendContainerState(codecSelfer_containerArrayEnd1819)
		return
	}
	z.DecSendContainerState(codecSelfer_containerArrayElem1819)
	if r.TryDecodeAsNil() {
		x.CreatedIndex = 0
	} else {
		yyv28 := &x.CreatedIndex
		yym29 := z.DecBinary()
		_ = yym29
		if false {
		} else {
			*((*uint64)(yyv28)) = uint64(r.DecodeUint(64))
		}
	}
	yyj20++
	if yyhl20 {
		yyb20 = yyj20 > l
	} else {
		yyb20 = r.CheckBreak()
	}
	if yyb20 {
		z.DecSendContainerState(codecSelfer_containerArrayEnd1819)
		return
	}
	z.DecSendContainerState(codecSelfer_containerArrayElem1819)
	if r.TryDecodeAsNil() {
		x.ModifiedIndex = 0
	} else {
		yyv30 := &x.ModifiedIndex
		yym31 := z.DecBinary()
		_ = yym31
		if false {
		} else {
			*((*uint64)(yyv30)) = uint64(r.DecodeUint(64))
		}
	}
	yyj20++
	if yyhl20 {
		yyb20 = yyj20 > l
	} else {
		yyb20 = r.CheckBreak()
	}
	if yyb20 {
		z.DecSendContainerState(codecSelfer_containerArrayEnd1819)
		return
	}
	z.DecSendContainerState(codecSelfer_containerArrayElem1819)
	if r.TryDecodeAsNil() {
		if x.Expiration != nil {
			x.Expiration = nil
		}
	} else {
		if x.Expiration == nil {
			x.Expiration = new(time.Time)
		}
		yym33 := z.DecBinary()
		_ = yym33
		if false {
		} else if yym34 := z.TimeRtidIfBinc(); yym34 != 0 {
			r.DecodeBuiltin(yym34, x.Expiration)
		} else if z.HasExtensions() && z.DecExt(x.Expiration) {
		} else if yym33 {
			z.DecBinaryUnmarshal(x.Expiration)
		} else if !yym33 && z.IsJSONHandle() {
			z.DecJSONUnmarshal(x.Expiration)
		} else {
			z.DecFallback(x.Expiration, false)
		}
	}
	yyj20++
	if yyhl20 {
		yyb20 = yyj20 > l
	} else {
		yyb20 = r.CheckBreak()
	}
	if yyb20 {
		z.DecSendContainerState(codecSelfer_containerArrayEnd1819)
		return
	}
	z.DecSendContainerState(codecSelfer_containerArrayElem1819)
	if r.TryDecodeAsNil() {
		x.TTL = 0
	} else {
		yyv35 := &x.TTL
		yym36 := z.DecBinary()
		_ = yym36
		if false {
		} else {
			*((*int64)(yyv35)) = int64(r.DecodeInt(64))
		}
	}
	for {
		yyj20++
		if yyhl20 {
			yyb20 = yyj20 > l
		} else {
			yyb20 = r.CheckBreak()
		}
		if yyb20 {
			break
		}
		z.DecSendContainerState(codecSelfer_containerArrayElem1819)
		z.DecStructFieldNotFound(yyj20-1, "")
	}
	z.DecSendContainerState(codecSelfer_containerArrayEnd1819)
}

func (x Nodes) CodecEncodeSelf(e *codec1978.Encoder) {
	var h codecSelfer1819
	z, r := codec1978.GenHelperEncoder(e)
	_, _, _ = h, z, r
	if x == nil {
		r.EncodeNil()
	} else {
		yym1 := z.EncBinary()
		_ = yym1
		if false {
		} else if z.HasExtensions() && z.EncExt(x) {
		} else {
			h.encNodes((Nodes)(x), e)
		}
	}
}

func (x *Nodes) CodecDecodeSelf(d *codec1978.Decoder) {
	var h codecSelfer1819
	z, r := codec1978.GenHelperDecoder(d)
	_, _, _ = h, z, r
	yym1 := z.DecBinary()
	_ = yym1
	if false {
	} else if z.HasExtensions() && z.DecExt(x) {
	} else {
		h.decNodes((*Nodes)(x), d)
	}
}

func (x codecSelfer1819) encNodes(v Nodes, e *codec1978.Encoder) {
	var h codecSelfer1819
	z, r := codec1978.GenHelperEncoder(e)
	_, _, _ = h, z, r
	r.EncodeArrayStart(len(v))
	for _, yyv1 := range v {
		z.EncSendContainerState(codecSelfer_containerArrayElem1819)
		if yyv1 == nil {
			r.EncodeNil()
		} else {
			yyv1.CodecEncodeSelf(e)
		}
	}
	z.EncSendContainerState(codecSelfer_containerArrayEnd1819)
}

func (x codecSelfer1819) decNodes(v *Nodes, d *codec1978.Decoder) {
	var h codecSelfer1819
	z, r := codec1978.GenHelperDecoder(d)
	_, _, _ = h, z, r

	yyv1 := *v
	yyh1, yyl1 := z.DecSliceHelperStart()
	var yyc1 bool
	_ = yyc1
	if yyl1 == 0 {
		if yyv1 == nil {
			yyv1 = []*Node{}
			yyc1 = true
		} else if len(yyv1) != 0 {
			yyv1 = yyv1[:0]
			yyc1 = true
		}
	} else if yyl1 > 0 {
		var yyrr1, yyrl1 int
		var yyrt1 bool
		_, _ = yyrl1, yyrt1
		yyrr1 = yyl1 // len(yyv1)
		if yyl1 > cap(yyv1) {

			yyrg1 := len(yyv1) > 0
			yyv21 := yyv1
			yyrl1, yyrt1 = z.DecInferLen(yyl1, z.DecBasicHandle().MaxInitLen, 8)
			if yyrt1 {
				if yyrl1 <= cap(yyv1) {
					yyv1 = yyv1[:yyrl1]
				} else {
					yyv1 = make([]*Node, yyrl1)
				}
			} else {
				yyv1 = make([]*Node, yyrl1)
			}
			yyc1 = true
			yyrr1 = len(yyv1)
			if yyrg1 {
				copy(yyv1, yyv21)
			}
		} else if yyl1 != len(yyv1) {
			yyv1 = yyv1[:yyl1]
			yyc1 = true
		}
		yyj1 := 0
		for ; yyj1 < yyrr1; yyj1++ {
			yyh1.ElemContainerState(yyj1)
			if r.TryDecodeAsNil() {
				if yyv1[yyj1] != nil {
					*yyv1[yyj1] = Node{}
				}
			} else {
				if yyv1[yyj1] == nil {
					yyv1[yyj1] = new(Node)
				}
				yyw2 := yyv1[yyj1]
				yyw2.CodecDecodeSelf(d)
			}

		}
		if yyrt1 {
			for ; yyj1 < yyl1; yyj1++ {
				yyv1 = append(yyv1, nil)
				yyh1.ElemContainerState(yyj1)
				if r.TryDecodeAsNil() {
					if yyv1[yyj1] != nil {
						*yyv1[yyj1] = Node{}
					}
				} else {
					if yyv1[yyj1] == nil {
						yyv1[yyj1] = new(Node)
					}
					yyw3 := yyv1[yyj1]
					yyw3.CodecDecodeSelf(d)
				}

			}
		}

	} else {
		yyj1 := 0
		for ; !r.CheckBreak(); yyj1++ {

			if yyj1 >= len(yyv1) {
				yyv1 = append(yyv1, nil) // var yyz1 *Node
				yyc1 = true
			}
			yyh1.ElemContainerState(yyj1)
			if yyj1 < len(yyv1) {
				if r.TryDecodeAsNil() {
					if yyv1[yyj1] != nil {
						*yyv1[yyj1] = Node{}
					}
				} else {
					if yyv1[yyj1] == nil {
						yyv1[yyj1] = new(Node)
					}
					yyw4 := yyv1[yyj1]
					yyw4.CodecDecodeSelf(d)
				}

			} else {
				z.DecSwallow()
			}

		}
		if yyj1 < len(yyv1) {
			yyv1 = yyv1[:yyj1]
			yyc1 = true
		} else if yyj1 == 0 && yyv1 == nil {
			yyv1 = []*Node{}
			yyc1 = true
		}
	}
	yyh1.End()
	if yyc1 {
		*v = yyv1
	}
}
