# go-xmpp

This library lets you open a XMPP stream over a TCP connection (or any other io.ReadWriter) and exchange XMPP protocol messages.

This is a work in progress in its early stages, some things still don't work, some might change drastically. Not for use in production environments.

## Quick start

```go
tcp, _ := net.Dial("tcp", "server.com:5222")

stream := &xmpp.StreamHeader{
    Namespace: xmpp.NamespaceClient,
    To:        "server.com",
}

conn, err := xmpp.Open(tcp, stream) // Open a bidirectional XMPP stream over the TCP connection (as a client)

conn.Features() // returns *xmpp.Features with stream features
conn.Write(&xmpp.StartTLS{}) // Send a starttls message
conn.Read() // Read an XMPP protocol message from the stream
conn.Close() // Close the XMPP stream

```