package com.fabric.listener;

import javax.servlet.http.HttpSession;
import javax.websocket.HandshakeResponse;
import javax.websocket.server.HandshakeRequest;
import javax.websocket.server.ServerEndpointConfig;
import javax.websocket.server.ServerEndpointConfig.Configurator;


/**
 * scoket获取session
 * @author Administrator
 *
 */
public class GetHttpSessionConfigurator extends Configurator  {

	@Override
	public void modifyHandshake(ServerEndpointConfig sec,
			HandshakeRequest request, HandshakeResponse response) {
		// TODO Auto-generated method stub
		 HttpSession httpSession=(HttpSession) request.getHttpSession();
	     sec.getUserProperties().put(HttpSession.class.getName(),httpSession);
	}
}