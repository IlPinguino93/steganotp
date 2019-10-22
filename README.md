# steganotp
This is a quick tool that combines One-Time-Pad (a symmetrical polialphabetic cipher from the cold war) with Steganography (The art of hiding things in images.)

## Differences to regular steganography
- No size differences, as the concealed data is not appended.
- No loss of quality on the conceald data
- The original image is required
- Without the original image, there is no chance of restoring the concealed data reliably - this is similiar to One time pad.

## Terminology
* Lock: Image with information encoded on it.
* Key: Image without information encoded on it. 
* Secret: The information encoded.
* Penguin: Flightless bird in Antarctica that has nothing to do with this project. Why do you even care?

## Regarding the "key" image
If you're old fashioned, you may want to think of the key image as a one-time pad. The size of the key image limits the secret size. One pixel equals one byte. If it's lost, you can't restore your secret. 

## How it works (theory)
### Encoding
1. You need a PNG image, this is your "key" image. Lose it, your data is unrecoverable.
    1. Each pixel of the image stores one byte of information. 
2. The byte is converted into a 3 digit octal number. This number is added or subtracted to/from the pixel color values(RGB). 
3. The algorithm produces a new image (the "Lock" image) that is encoded with the secret information. 

### Decoding
1. You need both the lock and the key image. They must be equal in size, obviously. 
2. Each pixel of the lock image is compared to the same pixel in the key image. The difference is one byte of your data. 
3. The data is decoded.

## Using this in your project

### Getting the image's possible area in bytes
```
steganotp.Size(key image.Image) int
```

### Encoding information into the image
```
steganotp.Encode(key image.Image, data byte[]) image.Image, err
```

### Decoding information using lock and key
```
steganotp.Decode(lock image.Image, key image.Image) byte[], err
```

## CLI

I'm going to build a simple CLI for this. I'll update this readme once it's done. 

## Disclaimer
This was built in one afternoon, mostly because I was bored and learning Golang. I'm not a particularcy good coder, and I'm not a genius that can replace a well-working, peer-reviewed encryption method. 

If you really need to encrypt something safely and reliably, rely on tested methods like RSA or AES. Don't use steganotp unless you're a curious enthusiast who likes the idea. 

