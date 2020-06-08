package com.fabric.model;

import java.io.Serializable;

/**
 * 
 * @author haihua
 * 
 *         用户模型
 * 
 *         窗口 > 首选项 > Java > 代码生成 > 代码和注释
 */

public class Message implements Serializable

{

	private static final long serialVersionUID = -4360427971861239742L;
	private String senderID;
	private String receiverID;
	private String sourceID;
	private String destinationID;
	private String message;
	private String file;
	private int msgNumber;
	
	
	
	public int getMsgNumber() {
		return msgNumber;
	}

	public void setMsgNumber(int msgNumber) {
		this.msgNumber = msgNumber;
	}
	
	public String getSourceID() {
		return sourceID;
	}

	public void setSourceID(String sourceID) {
		this.sourceID = sourceID;
	}

	public String getDestinationID() {
		return destinationID;
	}

	public void setDestinationID(String destinationID) {
		this.destinationID = destinationID;
	}

	public String getMessage() {
		return message;
	}

	public void setMessage(String message) {
		this.message = message;
	}

	public String getSenderID() {
		return senderID;
	}

	public void setSenderID(String senderID) {
		this.senderID = senderID;
	}

	public String getReceiverID() {
		return receiverID;
	}

	public void setReceiverID(String receiverID) {
		this.receiverID = receiverID;
	}

	public String getFile() {
		return file;
	}

	public void setFile(String file) {
		this.file = file;
	}
}