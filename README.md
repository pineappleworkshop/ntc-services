We need to implement something like this:
https://unisat.io/wallet-api-v4/address/inscriptions
?address=bc1p5pvvfjtnhl32llttswchrtyd9mdzd3p7yps98tlydh2dm6zj6gqsfkmcnd&cursor=0&size=10

We need to essentially dump redb to mongodb.

```json
{
  "list": [
    {
      "inscriptionId": "4e80d14abdb35ce193758cfd69ae8ce67f8036368ac75b729ef2fd3e0c6bad2fi0",
      "inscriptionNumber": 11822188,
      "address": "bc1p5pvvfjtnhl32llttswchrtyd9mdzd3p7yps98tlydh2dm6zj6gqsfkmcnd",
      "outputValue": 546,
      "preview": "https://ordinals.com/preview/4e80d14abdb35ce193758cfd69ae8ce67f8036368ac75b729ef2fd3e0c6bad2fi0",
      "content": "https://ordinals.com/content/4e80d14abdb35ce193758cfd69ae8ce67f8036368ac75b729ef2fd3e0c6bad2fi0",
      "contentLength": 57,
      "contentType": "text/plain;charset=utf-8",
      "contentBody": "",
      "timestamp": 1686720452,
      "genesisTransaction": "4e80d14abdb35ce193758cfd69ae8ce67f8036368ac75b729ef2fd3e0c6bad2f",
      "location": "5ea97a576a9a4368ec6f45d15cb5a1a4d65d68e060aea0c8b5be477e4ec59ea1:8:0",
      "output": "5ea97a576a9a4368ec6f45d15cb5a1a4d65d68e060aea0c8b5be477e4ec59ea1:8",
      "offset": 0
    },
    {
      "inscriptionId": "c0e627075a7991e5049230c886e841c9eb82f2b1dc392ec86acd706c25d72afdi0",
      "inscriptionNumber": 14938723,
      "address": "bc1p5pvvfjtnhl32llttswchrtyd9mdzd3p7yps98tlydh2dm6zj6gqsfkmcnd",
      "outputValue": 546,
      "preview": "https://ordinals.com/preview/c0e627075a7991e5049230c886e841c9eb82f2b1dc392ec86acd706c25d72afdi0",
      "content": "https://ordinals.com/content/c0e627075a7991e5049230c886e841c9eb82f2b1dc392ec86acd706c25d72afdi0",
      "contentLength": 4,
      "contentType": "text/plain;charset=utf-8",
      "contentBody": "",
      "timestamp": 1688653182,
      "genesisTransaction": "c0e627075a7991e5049230c886e841c9eb82f2b1dc392ec86acd706c25d72afd",
      "location": "5ea97a576a9a4368ec6f45d15cb5a1a4d65d68e060aea0c8b5be477e4ec59ea1:5:0",
      "output": "5ea97a576a9a4368ec6f45d15cb5a1a4d65d68e060aea0c8b5be477e4ec59ea1:5",
      "offset": 0
    },
    {
      "inscriptionId": "720dff8e5224c5f974918cab07fcbe1d820f6cca2cbc3d5bfff22cd7eb76eb7ci0",
      "inscriptionNumber": 14938719,
      "address": "bc1p5pvvfjtnhl32llttswchrtyd9mdzd3p7yps98tlydh2dm6zj6gqsfkmcnd",
      "outputValue": 546,
      "preview": "https://ordinals.com/preview/720dff8e5224c5f974918cab07fcbe1d820f6cca2cbc3d5bfff22cd7eb76eb7ci0",
      "content": "https://ordinals.com/content/720dff8e5224c5f974918cab07fcbe1d820f6cca2cbc3d5bfff22cd7eb76eb7ci0",
      "contentLength": 18,
      "contentType": "text/plain;charset=utf-8",
      "contentBody": "",
      "timestamp": 1688653182,
      "genesisTransaction": "720dff8e5224c5f974918cab07fcbe1d820f6cca2cbc3d5bfff22cd7eb76eb7c",
      "location": "5ea97a576a9a4368ec6f45d15cb5a1a4d65d68e060aea0c8b5be477e4ec59ea1:6:0",
      "output": "5ea97a576a9a4368ec6f45d15cb5a1a4d65d68e060aea0c8b5be477e4ec59ea1:6",
      "offset": 0
    },
    {
      "inscriptionId": "88665e98d24676cb2268551bda756dbfe79c0bb0706812fc4c6ebb5cdf31cf1ai0",
      "inscriptionNumber": 14938714,
      "address": "bc1p5pvvfjtnhl32llttswchrtyd9mdzd3p7yps98tlydh2dm6zj6gqsfkmcnd",
      "outputValue": 546,
      "preview": "https://ordinals.com/preview/88665e98d24676cb2268551bda756dbfe79c0bb0706812fc4c6ebb5cdf31cf1ai0",
      "content": "https://ordinals.com/content/88665e98d24676cb2268551bda756dbfe79c0bb0706812fc4c6ebb5cdf31cf1ai0",
      "contentLength": 36,
      "contentType": "text/plain;charset=utf-8",
      "contentBody": "",
      "timestamp": 1688653182,
      "genesisTransaction": "88665e98d24676cb2268551bda756dbfe79c0bb0706812fc4c6ebb5cdf31cf1a",
      "location": "5ea97a576a9a4368ec6f45d15cb5a1a4d65d68e060aea0c8b5be477e4ec59ea1:7:0",
      "output": "5ea97a576a9a4368ec6f45d15cb5a1a4d65d68e060aea0c8b5be477e4ec59ea1:7",
      "offset": 0
    }
  ],
  "total": 4
}
```