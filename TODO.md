- More efficient file handling, see TODOs in code

## Security

- secure the connection

## CI

- make file exacutable
- statically link?

## Maybe

- If multiple files are received send each one to a different fd

## Network

- Send function will send arbitrary amount of data, no matter the chunk size
  - It will work by sending data through a channel in the stream struct
  - File/Stdin sending will send in chunk size or less always
    - Every chunk will have a header and is encrypted
  - While initiliazing, the data will be sent with header but not encrypted -- flag in the stream to turn on encryption
  - A key will be in the stream struct used for symmetric encryption
- Receive functions
  - One function/variant will receive all the data at once becuase its pointless receiveing by chunks if the data is some signature/key
  - Other function/variant will receive data and send it to a channel which is consumed by a printer/writer
    - This functions needs to decrypt the data
    - This functions needs handle a situation where the header describes data longer the a chunk, and where the chunk is imidietly follwed by another chunk with a header
      - At first the initial chunk needs to be processed and then the rest of the read data is counted as a new chunk
