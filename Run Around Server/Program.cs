using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Net;
using System.Net.Sockets;
using System.IO;
using System.Threading;
namespace Run_Around_Server
{
    class Program
    {
        static void Main(string[] args)
        {

            IPEndPoint UplinkEP = new IPEndPoint(IPAddress.Any, 8443);
            IPEndPoint DownLinkEP = new IPEndPoint(IPAddress.Any, 8880);
            Socket UplinkSock = new Socket(AddressFamily.InterNetwork, SocketType.Stream, ProtocolType.Tcp);
            Socket DownLinkSock = new Socket(AddressFamily.InterNetwork, SocketType.Stream, ProtocolType.Tcp);
            UplinkSock.Bind(UplinkEP);
            UplinkSock.Listen(10);
            DownLinkSock.Bind(DownLinkEP);
            DownLinkSock.Listen(10);
            
            // Both ends setup.

            Socket DownlinkSockR = DownLinkSock.Accept();
            Console.WriteLine("Downlink online.");
            Socket UplinkSockR = UplinkSock.Accept();
            Console.WriteLine("Uplink online. Making Connection to Final Endpoint.");

            Socket ProxyEndSock = new Socket(AddressFamily.InterNetwork, SocketType.Stream, ProtocolType.Tcp);
            IPEndPoint ProxyEndEP = new IPEndPoint(IPAddress.Parse("127.0.0.1"), 22);
            /*
            NetworkStream DownLinkNS = new NetworkStream(DownlinkSockR);
            NetworkStream UpLinkNS = new NetworkStream(UplinkSockR);
            StreamReader DownLinkSr = new StreamReader(DownLinkNS);
            StreamWriter DownLinkSw = new StreamWriter(DownLinkNS);
            StreamReader UpLinkSr = new StreamReader(UpLinkNS);*/
            Console.WriteLine("Made Streams");
            while (true)
            {
                int DLSA = DownlinkSockR.Available;
                if (DLSA != 0){
                    Thread.Sleep(33);
                    Console.WriteLine("Headers are {0} bytes long", DLSA);
                    byte[] Inbound = new byte[DLSA];
                    DownlinkSockR.Receive(Inbound);
                    break;
                }
            }
            Console.WriteLine("Got the Clients (bullshit) headers");
            string ResponceHeaders = @"HTTP/1.0 200 OK
Vary: Accept-Encoding
Content-Type: application/octet-stream
Accept-Ranges: bytes
Last-Modified: Tue, 11 Dec 2012 10:32:35 GMT
Content-Length: 1073741824
Date: Wed, 30 Jan 2013 10:59:29 GMT
Server: SeemsLegitCo/1
Proxy-Connection: keep-alive

";
            DownlinkSockR.Send(Encoding.ASCII.GetBytes(ResponceHeaders));
            Console.WriteLine("Sent my headers");
            // OKay So the proxy thinks we are legit as fuck now.
            // Lets connect to the ProxyEndpoint.
            ProxyEndSock.Connect(ProxyEndEP);
            Console.WriteLine("We have connected to our localend.");
            NetworkStream ProxyEndNS = new NetworkStream(ProxyEndSock);
            StreamReader ProxyEndSR = new StreamReader(ProxyEndNS);
            StreamWriter ProxyEndSW = new StreamWriter(ProxyEndNS);

            while (true)
            {
                Thread.Sleep(33); // hacks hacks lazyness etc
                int UplinkAvail = UplinkSockR.Available;
                if (UplinkAvail != 0)
                {
                    // read and send
                    byte[] Inbound = new byte[UplinkAvail];
                    UplinkSockR.Receive(Inbound);
                    ProxyEndSock.Send(Inbound);
                }
                int ProxyAvail = ProxyEndSock.Available;
                if (ProxyAvail != 0)
                {
                    // read and send.
                    byte[] Inbound = new byte[ProxyAvail];
                    ProxyEndSock.Receive(Inbound);
                    DownlinkSockR.Send(Inbound);
                }
            }
            
        }
    }
}
