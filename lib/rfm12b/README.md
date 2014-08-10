rfm12b Driver for go
====================

Simple serial driver for RFM12b sketch in GO

requires logger library

provides struct {NodeId, MsgData[]bytes} via channel provided in start.

Example

```


const (
  //portName string = "/dev/ttyUSB0" //rPi USB
  //portName string = "/dev/ttyAMA0" //rPi Header
  portName string = "/dev/tty.usbserial-A1014KGL" //Mac
  baud uint32 = 57600
)

func main() {
  jeelink := rfm12b.New(portName, baud, "./Logs/RFM12b/")

  cs := make(chan rfm12b.Output)

  go jeelink.Start(cs)

  for m := range cs{
        fmt.Println(m.NodeId, m.MsgData)
  }

}
```
