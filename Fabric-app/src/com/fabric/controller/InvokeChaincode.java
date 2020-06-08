/****************************************************** 
 *  Copyright 2018 IBM Corporation 
 *  Licensed under the Apache License, Version 2.0 (the "License"); 
 *  you may not use this file except in compliance with the License. 
 *  You may obtain a copy of the License at 
 *  http://www.apache.org/licenses/LICENSE-2.0 
 *  Unless required by applicable law or agreed to in writing, software 
 *  distributed under the License is distributed on an "AS IS" BASIS, 
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. 
 *  See the License for the specific language governing permissions and 
 *  limitations under the License.
 */
package com.fabric.controller;

import static java.nio.charset.StandardCharsets.UTF_8;

import java.util.Collection;
import java.util.HashMap;
import java.util.Map;
import java.util.Vector;
import java.util.logging.Level;
import java.util.logging.Logger;
import java.util.regex.Pattern;




import main.java.org.example.client.CAClient;
import main.java.org.example.client.ChannelClient;
import main.java.org.example.client.FabricClient;
import main.java.org.example.config.Config;
import main.java.org.example.user.UserContext;
import main.java.org.example.util.Util;
import net.sf.json.JSONObject;

import org.hyperledger.fabric.sdk.ChaincodeEventListener;
import org.hyperledger.fabric.sdk.ChaincodeID;
import org.hyperledger.fabric.sdk.ChaincodeResponse.Status;
import org.hyperledger.fabric.sdk.Channel;
import org.hyperledger.fabric.sdk.EventHub;
import org.hyperledger.fabric.sdk.Orderer;
import org.hyperledger.fabric.sdk.Peer;
import org.hyperledger.fabric.sdk.ProposalResponse;
import org.hyperledger.fabric.sdk.TransactionProposalRequest;
import org.hyperledger.fabric.sdk.BlockEvent;
import org.hyperledger.fabric.sdk.ChaincodeEvent;

import com.fabric.model.ChaincodeEventCapture;
import com.fabric.model.Message;

public class InvokeChaincode {

	private static final byte[] EXPECTED_EVENT_DATA = "!".getBytes(UTF_8);
	private static final String EXPECTED_EVENT_NAME = "msgEvent";

	public static String Invoke_createMsg(Message msg, String function) {
		try {
			Util.cleanUp();
			String caUrl = Config.CA_ORG1_URL;
			CAClient caClient = new CAClient(caUrl, null);
			// Enroll Admin to Org1MSP
			UserContext adminUserContext = new UserContext();
			adminUserContext.setName(Config.ADMIN);
			adminUserContext.setAffiliation(Config.ORG1);
			adminUserContext.setMspId(Config.ORG1_MSP);
			caClient.setAdminUserContext(adminUserContext);
			adminUserContext = caClient.enrollAdminUser(Config.ADMIN, Config.ADMIN_PW);

			FabricClient fabClient = new FabricClient(adminUserContext);

			ChannelClient channelClient = fabClient.createChannelClient(Config.CHANNEL_NAME);
			Channel channel = channelClient.getChannel();
			Peer peer = fabClient.getInstance().newPeer(Config.ORG1_PEER_0, Config.ORG1_PEER_0_URL);
			EventHub eventHub = fabClient.getInstance().newEventHub("eventhub01", "grpc://localhost:7053");
			Orderer orderer = fabClient.getInstance().newOrderer(Config.ORDERER_NAME, Config.ORDERER_URL);
			channel.addPeer(peer);
			channel.addEventHub(eventHub);
			channel.addOrderer(orderer);
			channel.initialize();
			String event = "";
			final Vector<ChaincodeEventCapture> chaincodeEvents = new Vector<>(); // Test
																					// list
																					// to
																					// capture

			ChaincodeEventListener chaincodeEventListener = new ChaincodeEventListener() {

				@Override
				public void received(String handle, BlockEvent blockEvent, ChaincodeEvent chaincodeEvent) {
					chaincodeEvents.add(new ChaincodeEventCapture(handle, blockEvent, chaincodeEvent));
					// System.out.println(blockEvent);
				}
			};
			
			// chaincode events.
			String eventListenerHandle = channel.registerChaincodeEventListener(Pattern.compile(".*"), 
					Pattern.compile(Pattern.quote(EXPECTED_EVENT_NAME)),chaincodeEventListener);

			TransactionProposalRequest request = fabClient.getInstance().newTransactionProposalRequest();
			ChaincodeID ccid = ChaincodeID.newBuilder().setName(Config.CHAINCODE_1_NAME).build();
			request.setChaincodeID(ccid);
			request.setFcn(function);
			// String[] arguments = { msg.getMessageID(), msg.getName(),
			// msg.getContent(), msg.getNetworkID(), msg.getTimestamp() };

			String[] arguments = { msg.getSenderID(), msg.getReceiverID(),
					msg.getSourceID(), msg.getDestinationID(),
					msg.getMessage(), msg.getFile(), Integer.toString(msg.getMsgNumber())};

			// System.out.println(msg.getSenderID() + msg.getReceiverID() +
			// msg.getSourceID() + msg.getDestinationID() + msg.getMessage());

			request.setArgs(arguments);
			// request.setProposalWaitTime(1000);

			Map<String, byte[]> tm2 = new HashMap<String, byte[]>();
			tm2.put("HyperLedgerFabric", "TransactionProposalRequest:JavaSDK".getBytes(UTF_8));
			tm2.put("method", "TransactionProposalRequest".getBytes(UTF_8));
			tm2.put("result", ":)".getBytes(UTF_8));
			tm2.put(EXPECTED_EVENT_NAME, EXPECTED_EVENT_DATA);
			request.setTransientMap(tm2);
			Collection<ProposalResponse> responses = channelClient.sendTransactionProposal(request);

			for (ProposalResponse res : responses) {
				Status status = res.getStatus();
				Logger.getLogger(InvokeChaincode.class.getName()).log(
						Level.INFO,"Invoked createMessage on " + Config.CHAINCODE_1_NAME + ". Status - " + status);
			}
			Thread.sleep(5000);
			if (eventListenerHandle != null) {
				for (ChaincodeEventCapture chaincodeEventCapture : chaincodeEvents) {
					String payload = new String(chaincodeEventCapture.getChaincodeEvent().getPayload());
					JSONObject jsonEvent = JSONObject.fromObject(payload);
					JSONObject jsonEventPlus = new JSONObject();
					jsonEventPlus.put("event", jsonEvent);
					//System.out.println(jsonEventPlus);
					event = jsonEventPlus.toString();
				}
				chaincodeEvents.clear();
			}
			return event;

		} catch (Exception e) {
			e.printStackTrace();
		}
		return null;
	}
	
	
	public static String Invoke_receiveMsg(String encEvent, String function) {
		try {
			Util.cleanUp();
			String caUrl = Config.CA_ORG1_URL;
			CAClient caClient = new CAClient(caUrl, null);
			// Enroll Admin to Org1MSP
			UserContext adminUserContext = new UserContext();
			adminUserContext.setName(Config.ADMIN);
			adminUserContext.setAffiliation(Config.ORG1);
			adminUserContext.setMspId(Config.ORG1_MSP);
			caClient.setAdminUserContext(adminUserContext);
			adminUserContext = caClient.enrollAdminUser(Config.ADMIN, Config.ADMIN_PW);

			FabricClient fabClient = new FabricClient(adminUserContext);

			ChannelClient channelClient = fabClient.createChannelClient(Config.CHANNEL_NAME);
			Channel channel = channelClient.getChannel();
			Peer peer = fabClient.getInstance().newPeer(Config.ORG1_PEER_0, Config.ORG1_PEER_0_URL);
			EventHub eventHub = fabClient.getInstance().newEventHub("eventhub01", "grpc://localhost:7053");
			Orderer orderer = fabClient.getInstance().newOrderer(Config.ORDERER_NAME, Config.ORDERER_URL);
			channel.addPeer(peer);
			channel.addEventHub(eventHub);
			channel.addOrderer(orderer);
			channel.initialize();
			String event = "";
			final Vector<ChaincodeEventCapture> chaincodeEvents = new Vector<>(); // Test
																					// list
																					// to
																					// capture

			ChaincodeEventListener chaincodeEventListener = new ChaincodeEventListener() {

				@Override
				public void received(String handle, BlockEvent blockEvent, ChaincodeEvent chaincodeEvent) {
					chaincodeEvents.add(new ChaincodeEventCapture(handle, blockEvent, chaincodeEvent));
					// System.out.println(blockEvent);
				}
			};
			
			// chaincode events.
			String eventListenerHandle = channel.registerChaincodeEventListener(Pattern.compile(".*"), 
					Pattern.compile(Pattern.quote("receiveEvent")),chaincodeEventListener);

			TransactionProposalRequest request = fabClient.getInstance().newTransactionProposalRequest();
			ChaincodeID ccid = ChaincodeID.newBuilder().setName(Config.CHAINCODE_1_NAME).build();
			request.setChaincodeID(ccid);
			request.setFcn(function);
			
			JSONObject jsonEvent = JSONObject.fromObject(encEvent);
			JSONObject jsonEventPlus = jsonEvent.getJSONObject("event");
			System.out.println(jsonEventPlus);
			request.setArgs(jsonEventPlus.toString(), Integer.toString(MessageController.msgNumber));

			Map<String, byte[]> tm2 = new HashMap<String, byte[]>();
			tm2.put("HyperLedgerFabric", "TransactionProposalRequest:JavaSDK".getBytes(UTF_8));
			tm2.put("method", "TransactionProposalRequest".getBytes(UTF_8));
			tm2.put("result", ":)".getBytes(UTF_8));
			tm2.put("receiveEvent", EXPECTED_EVENT_DATA);
			request.setTransientMap(tm2);
			Collection<ProposalResponse> responses = channelClient.sendTransactionProposal(request);

			for (ProposalResponse res : responses) {
				Status status = res.getStatus();
				Logger.getLogger(InvokeChaincode.class.getName()).log(
						Level.INFO,"Invoked createMessage on " + Config.CHAINCODE_1_NAME + ". Status - " + status);
			}
			Thread.sleep(5000);
			if (eventListenerHandle != null) {
				for (ChaincodeEventCapture chaincodeEventCapture : chaincodeEvents) {
					String payload = new String(chaincodeEventCapture.getChaincodeEvent().getPayload(), "utf-8");
					event = payload;
				}
				chaincodeEvents.clear();
			}
			return event;

		} catch (Exception e) {
			e.printStackTrace();
		}
		return null;
	}
}
