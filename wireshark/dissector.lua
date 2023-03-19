-- declare our protocol
local pqcom_proto = Proto("PQCOM", "Post-Quantum Communication Protocol")
TypeTbl = {
  [0] = "ClientInitT",
  [1] = "ServerInitT",
  [2] = "ContentT",
  [3] = "ErrorT",
}
KemsNameTbl = {
  [0] = "PqComKyber512",
  [1] = "CirclKyber512",
  [2] = "CirclKyber768",
  [3] = "CirclKyber1024",
}
SignsNameTbl = {
  [0] = "PqComDilithium2",
  [1] = "CirclDilithium2",
  [2] = "CirclDilithium3",
  [3] = "CirclDilithium5",
}
KemsEkLenTbl = {
  [0] = 800,
  [1] = 800,
  [2] = 1184,
  [3] = 1568,
}
SignsSigLenTbl = {
  [0] = 2420,
  [1] = 2420,
  [2] = 3293,
  [3] = 4595,
}
KeyCiphertextTbl = {
  [0] = 768,
  [1] = 768,
  [2] = 1088,
  [3] = 1568,
}
NonceLen = 12

KemType = -1
SignType = -1

-- create a function to dissect it
function pqcom_proto.dissector(buffer, pinfo, tree)
  pinfo.cols.protocol = "PQCOM"
  local subtree = tree:add(pqcom_proto, buffer(), "Post-Quantum Communication Protocol")
  local type = buffer(2, 1):uint()
  local len = buffer(0, 2):uint()
  subtree:add(buffer(0, 2), "Length: " .. len)
  subtree:add(buffer(2, 1), string.format("Type: %s (%d)", TypeTbl[type], type))

  -- ClientInitT
  if type == 0 then
    -- KEM Type
    KemType = buffer(3, 1):uint()
    subtree:add(buffer(3, 1), string.format("KEM Type: %s (%d)", KemsNameTbl[KemType], KemType))

    -- Sign Type
    SignType = buffer(4, 1):uint()
    subtree:add(buffer(4, 1), string.format("Signature Type: %s (%d)", SignsNameTbl[SignType], SignType))

    -- Timestamp
    local timestamp = buffer(5, 8):uint64()
    local date = os.date(_, tonumber(tostring(timestamp):sub(1, 10)))
    subtree:add(buffer(5, 8), "Timestamp: " .. date .. " (" .. timestamp .. ")")

    -- Ek
    local ekLen = KemsEkLenTbl[KemType]
    local offset = ekLen + 13
    subtree:add(buffer(13, ekLen), "Public Encryption Key: " .. buffer(13, ekLen))

    -- Nonce
    subtree:add(buffer(offset, NonceLen), "Nonce: " .. buffer(offset, NonceLen))
    offset = offset + NonceLen

    -- Signature
    local sigLen = SignsSigLenTbl[SignType]
    subtree:add(buffer(offset, sigLen), "Signature: " .. buffer(offset, sigLen))

    -- ServerInitT
  elseif type == 1 then
    -- Key Ciphertext
    local offset = 3
    local ctLen = KeyCiphertextTbl[KemType]
    subtree:add(buffer(offset, ctLen), "Key ciphertext: " .. buffer(offset, ctLen))
    offset = offset + ctLen

    -- Signature
    local sigLen = SignsSigLenTbl[SignType]
    subtree:add(buffer(offset, sigLen), "Signature: " .. buffer(offset, sigLen))

    -- ErrorT
  elseif type == 3 then
    subtree:add(buffer(3, len), "Error reason: " .. buffer(3, len):string())

    -- ContentT
  else
    subtree:add(buffer(3, len), "Data: " .. buffer(3, len))
  end
end

-- load the udp.port table
local tcp_table = DissectorTable.get("tcp.port")
-- register our protocol to handle tcp port 4444
tcp_table:add(4040, pqcom_proto)
