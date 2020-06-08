package com.fabric.controller;

import java.net.InetSocketAddress;
import java.nio.ByteBuffer;
import java.util.HashMap;
import java.util.Map;

import main.java.org.example.config.Config;
import net.sf.json.JSONObject;

import org.java_websocket.WebSocket;
import org.java_websocket.handshake.ClientHandshake;
import org.java_websocket.server.WebSocketServer;

public class SimpleServer extends WebSocketServer {
	
//	static int msgNumber = 0;
	public Map<String, WebSocket> connPool = new HashMap<String, WebSocket>();

	public SimpleServer(InetSocketAddress address) {
		super(address);
	}
		

	@Override
	public void onOpen(WebSocket conn, ClientHandshake handshake) {
		//conn.send("Welcome to the server!"); //This method sends a message to the new client
		//broadcast( "new connection: " + handshake.getResourceDescriptor() ); //This method sends a message to all clients connected
		System.out.println("new connection to " + conn.getRemoteSocketAddress().toString());
		connPool.put(conn.getRemoteSocketAddress().getHostString(), conn);
	}

	@Override
	public void onClose(WebSocket conn, int code, String reason, boolean remote) {
		System.out.println("closed " + conn.getRemoteSocketAddress() + " with exit code " + code + " additional info: " + reason);
		connPool.remove(conn.getRemoteSocketAddress().getHostString());
	}

	@Override
	public void onMessage(WebSocket conn, String message) {
		JSONObject jsonEvent = JSONObject.fromObject(message);
		JSONObject jsonEventPlus = jsonEvent.getJSONObject("event");
		System.out.println(jsonEventPlus);
		
		//System.out.println("received a encrypted message from bridge: "	+ conn.getRemoteSocketAddress() + ": " + message);
		try {
			MessageController.receiveMessage(message);
		} catch (Exception e) {
			// TODO Auto-generated catch block
			e.printStackTrace();
		}
		
	}

	@Override
	public void onMessage( WebSocket conn, ByteBuffer file ) {
		System.out.println("received ByteBuffer from "	+ conn.getRemoteSocketAddress());
	}

	@Override
	public void onError(WebSocket conn, Exception ex) {
		System.err.println("an error occurred on connection " + conn.getRemoteSocketAddress()  + ":" + ex);
	}
	
	@Override
	public void onStart() {
		System.out.println("server started successfully");
	}

	public static void main(String[] args) {
		String host = Config.LocalAddress;
		int port = 8081;

		WebSocketServer server = new SimpleServer(new InetSocketAddress(host, port));
		server.run();
	}
}