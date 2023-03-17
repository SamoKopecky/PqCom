-- declare our protocol
pqcom_proto = Proto("PqCom", "Post-Qunatum Communication Protocol")
-- create a function to dissect it
function pqcom_proto.dissector(buffer, pinfo, tree)
  local type_tbl = {
    [0] = "ClientInitT",
    [1] = "ServerInitT",
    [2] = "ContentT",
    [3] = "ErrorT",
  }
  pinfo.cols.protocol = "PqCom"
  local subtree = tree:add(pqcom_proto, buffer(), "Post-Qunatum Communication Protocol")
  local type = buffer(2, 1):uint()
  subtree:add(buffer(0, 2), "Length: " .. buffer(0, 2):uint())
  subtree:add(buffer(2, 1), "Type: " .. type_tbl[type] .. " (" .. type .. ")")
end

-- load the udp.port table
tcp_table = DissectorTable.get("tcp.port")
-- register our protocol to handle tcp port 4444
tcp_table:add(4040, pqcom_proto)
