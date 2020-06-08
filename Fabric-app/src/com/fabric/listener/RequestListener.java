package com.fabric.listener;

import java.net.InetSocketAddress;
import java.util.Collection;
import java.util.Iterator;
import java.util.Map;
import java.util.Map.Entry;
import java.util.Vector;

import javax.servlet.ServletContextEvent;
import javax.servlet.ServletContextListener;
import javax.servlet.annotation.WebListener;

import main.java.org.example.config.Config;

import org.java_websocket.WebSocket;

import com.fabric.controller.SimpleServer;

@WebListener
public class RequestListener implements ServletContextListener {
	static SimpleServer WSserver;
	@Override
	public void contextDestroyed(ServletContextEvent arg0) {
		// TODO Auto-generated method stub

	}

	@Override
	public void contextInitialized(ServletContextEvent servletContextEvent) {
		WSserver = new SimpleServer(new InetSocketAddress(Config.LocalAddress, 8081));
		
		Thread WSThread = new Thread(new Runnable(){

			@Override
			public void run() {
				// TODO Auto-generated method stub
				WSserver.run();
			}
			
		});
		WSThread.start();
	}
	
	public static void broadcastToBridge(String event){
		Collection<WebSocket> bridges = new Vector<>();
		Iterator<Entry<String, WebSocket>> conns = WSserver.connPool.entrySet().iterator();
		while (conns.hasNext()) {
			Map.Entry<String, WebSocket> entry = conns.next();
			if (Config.Bridge1.equals(entry.getKey())){
				bridges.add(entry.getValue());
			}
			else if (Config.Bridge2.equals(entry.getKey())){
				bridges.add(entry.getValue());
			}
			else{}
		}
		WSserver.broadcast(event, bridges);
	}
	
	public static void broadcastToFront(String event){
		Collection<WebSocket> frontEnd = new Vector<>();
		Iterator<Entry<String, WebSocket>> conns = WSserver.connPool.entrySet().iterator();
		while (conns.hasNext()) {
			Map.Entry<String, WebSocket> entry = conns.next();
			if (Config.Bridge1.equals(entry.getKey())){
			}
			else if (Config.Bridge2.equals(entry.getKey())){
			}
			else{
				frontEnd.add(entry.getValue());
			}
		}
		//System.out.println(frontEnd);
		WSserver.broadcast(event, frontEnd);
	}
	
}