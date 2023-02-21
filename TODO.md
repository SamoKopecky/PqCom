- More efficient file handling, see TODOs in code

## Security

- Check the type of the messages, if they are not correct throw fatal error
- Signature for init msgs
- Modularity of PQ algorithms
- Maybe Time/nonce at the init msgs

## CI

- make file exacutable
- statically link?

## Network

- Receive functions
  - One function/variant will receive all the data at once becuase its pointless receiveing by chunks if the data is some signature/key
  - Other function/variant will receive data and send it to a channel which is consumed by a printer/writer
    - This functions needs to decrypt the data
    - This functions needs handle a situation where the header describes data longer the a chunk, and where the chunk is imidietly follwed by another chunk with a header
      - At first the initial chunk needs to be processed and then the rest of the read data is counted as a new chunk
