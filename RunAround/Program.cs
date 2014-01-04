using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Net;
using System.Net.Sockets;
using System.IO;
using System.Threading;

namespace RunAround
{
    class Program
    {
        static void Main(string[] args)
        {

            string IP = "78.129.218.173";
            string UpLinkHeaders = @"GET / HTTP/1.1
Host: "+IP+@"
Connection: keep-alive
Accept: text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8
User-Agent: Mozilla/5.0 (Windows NT 6.2; WOW64) AppleWebKit/537.17 (KHTML, like Gecko) Chrome/24.0.1312.56 Safari/537.17
DNT: 1
Accept-Encoding: gzip,deflate,sdch
Accept-Language: en-US,en;q=0.8
Accept-Charset: ISO-8859-1,utf-8;q=0.7,*;q=0.3

";
            IPEndPoint ListenEndPoint = new IPEndPoint(IPAddress.Any,2222);
            Socket ListenSocket = new Socket(AddressFamily.InterNetwork,SocketType.Stream, ProtocolType.Tcp);
            ListenSocket.Bind(ListenEndPoint);
            ListenSocket.Listen(10);
            Console.WriteLine("Ready. Connect on {0}",ListenEndPoint.Port);
            Socket LocalClient = ListenSocket.Accept();

            // Now we have the local connection we need to knock at the other servers door.
            Socket RemoteDownLink = new Socket(AddressFamily.InterNetwork, SocketType.Stream, ProtocolType.Tcp);
            Socket RemoteUpLink = new Socket(AddressFamily.InterNetwork, SocketType.Stream, ProtocolType.Tcp);
            IPEndPoint ProxyDLEnd = new IPEndPoint(IPAddress.Parse(IP), 80);
            IPEndPoint ProxyULEnd = new IPEndPoint(IPAddress.Parse(IP), 443);

            try
            {
                RemoteDownLink.Connect(ProxyDLEnd);
                RemoteUpLink.Connect(ProxyULEnd);
                Console.WriteLine("Both Endpoints connected.\r\nSending Headers on HTTP end");
                RemoteDownLink.Send(Encoding.ASCII.GetBytes(UpLinkHeaders));
            }
            catch(Exception e)
            {
                Console.WriteLine(e.Message);
                Console.WriteLine(e.Source);
                Console.WriteLine(e.StackTrace);
                Console.ReadLine();
            }
            // Remote End bits
            /*NetworkStream DownLinkStream = new NetworkStream(RemoteDownLink);
            StreamReader DownLinkInboundStream = new StreamReader(DownLinkStream);
            NetworkStream UpLinkStream = new NetworkStream(RemoteUpLink);
            StreamWriter UpLinkWriteStream = new StreamWriter(UpLinkStream);*/
            // Local End Bits
            /*NetworkStream LocalEndStream = new NetworkStream(LocalClient);
            StreamReader LocalEndRead = new StreamReader(LocalEndStream);
            StreamWriter LocalEndWrite = new StreamWriter(LocalEndStream);*/

            // We need to wait for the headers that the downlink sends to finish.
            while (true)
            {
                Thread.Sleep(33); // ADSL anyway.
                int DLSA = RemoteDownLink.Available;
                if (DLSA != 0)
                {
                    Console.WriteLine("Headers are {0} bytes long", DLSA);
                    byte[] Inbound = new byte[DLSA];
                    RemoteDownLink.Receive(Inbound);
                    break;
                }
                //if (line == "") { break; } // HTTP sends a blank line at the end.
            }
            Console.WriteLine("Got Servers headers...");
            // Ok the headers are done.
            while (true)
            {
                Thread.Sleep(33); // ADSL anyway.
                int LESA = LocalClient.Available;
                if (LESA != 0)
                {
                    Console.WriteLine("-> {0}", LESA);
                    byte[] Inbound = new byte[LESA]; 
                    LocalClient.Receive(Inbound);
                    RemoteUpLink.Send(Inbound, LESA, SocketFlags.None);
                }
                int DLSA = RemoteDownLink.Available;
                if (DLSA != 0)
                {
                    Console.WriteLine("<- {0}", DLSA);
                    byte[] Inbound = new byte[DLSA];
                    RemoteDownLink.Receive(Inbound);
                    LocalClient.Send(Inbound, DLSA, SocketFlags.None);
                }
            }

        }
    }
}
