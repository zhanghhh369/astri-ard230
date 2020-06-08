package com.fabric.controller;


import java.util.Base64;

import net.sf.json.JSONObject;

import org.springframework.stereotype.Controller;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RequestMethod;
import org.springframework.web.bind.annotation.RequestParam;
import org.springframework.web.multipart.MultipartFile;
import org.springframework.web.servlet.ModelAndView;

import com.fabric.listener.RequestListener;
import com.fabric.model.Message;

@Controller
public class MessageController {

	// static String eventQueue = "";
	
	static int msgNumber = 0;

	@RequestMapping(value = "/index")
	public ModelAndView welcomePage() {
		ModelAndView index = new ModelAndView("index");
		return index;
	}

	
	@RequestMapping(value = "/send")
	public ModelAndView sendPage() {
		ModelAndView sendPage = new ModelAndView("sendMsg");
		return sendPage;
	}
	
	
	@RequestMapping(value = "/sendMessage", method = RequestMethod.POST)
	protected ModelAndView sendMessage(
			@RequestParam("senderID") String senderID,
			@RequestParam("receiverID") String receiverID,
			@RequestParam("sourceID") String sourceID,
			@RequestParam("destinationID") String destinationID,
			@RequestParam("message") String message,
			@RequestParam("file") MultipartFile file) throws Exception {
		Message msg = new Message();
		msg.setSenderID(senderID);
		msg.setReceiverID(receiverID);
		msg.setSourceID(sourceID);
		msg.setDestinationID(destinationID);
		msg.setMessage(message);
		msg.setMsgNumber(msgNumber);
		byte [] byteFile = file.getBytes();
		//RequestListener.broadcast(event + ";" + byteFile.toString());
		String stringFile = Base64.getEncoder().encodeToString(byteFile);
		msg.setFile(stringFile);
		String event = InvokeChaincode.Invoke_createMsg(msg, "createMsg");
		msgNumber ++;
		RequestListener.broadcastToBridge(event);
		ModelAndView Suc = new ModelAndView("success");
		Suc.addObject("event", event);
		return Suc;
	}

	@RequestMapping(value = "/observer")
	protected ModelAndView observer() throws Exception {
		System.out.println(QueryChaincode.Query(null));
		ModelAndView observer = new ModelAndView();
		observer.addObject("result", QueryChaincode.Query(null));
		observer.setViewName("observer");
		return observer;
	}
	
	protected static synchronized void receiveMessage(String encEvent) throws Exception {
	    
		String event = InvokeChaincode.Invoke_receiveMsg(encEvent, "receiveMsg");
		System.out.println(event);
		if(!"\"redundant\"".equals(event) && 
				!"\"Hash\"".equals(event) &&
				!"\"TransactionID\"".equals(event) &&
				!"\"SenderID\"".equals(event) &&
				!"\"ReceiverID\"".equals(event) &&
				!"\"Source\"".equals(event)&&
				!"\"DestinationID\"".equals(event)&&
				!"\"signature\"".equals(event)){
			JSONObject json_Dec_Event = JSONObject.fromObject(event);
			msgNumber = Integer.parseInt(json_Dec_Event.getJSONObject("header").getString("index")) + 1;
		}
		RequestListener.broadcastToFront(event);
		
	}

}
