package com.fabric.controller;

import java.net.InetSocketAddress;

import org.java_websocket.server.WebSocketServer;

public class WSHandler extends Thread {
	
	
	
	
	public WebSocketServer startServer() {
		String host = "localhost";
		int port = 8081;

		WebSocketServer server = new SimpleServer(new InetSocketAddress(host, port));
		server.run();
		return server;
}
}
