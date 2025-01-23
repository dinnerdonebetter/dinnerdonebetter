import * as crypto from 'crypto';

const encryptionAlgo = 'aes-256-cbc';
const keyLength = 32;

export class EncryptorDecryptor<T> {
  secretKey: string;
  initializationVectors: Buffer;

  constructor(encryptionKey?: string, initializationVectors?: string) {
    if ((encryptionKey || '').length !== keyLength) {
      console.log('encryption key not provided, generating random bytes');
      this.secretKey = crypto.randomBytes(keyLength).toString('hex').slice(0, keyLength);
    } else {
      this.secretKey = encryptionKey!;
    }

    if (Buffer.from(initializationVectors || '', 'base64').length !== keyLength) {
      console.log('initialization vectors not provided, generating random bytes');
      this.initializationVectors = crypto.randomBytes(keyLength);
    } else {
      this.initializationVectors = Buffer.from(initializationVectors!, 'base64');
    }
  }

  encrypt(x: T): string {
    let cipher = crypto.createCipheriv(
      encryptionAlgo,
      this.secretKey,
      this.initializationVectors.toString('hex').slice(0, 16),
    );

    const value = JSON.stringify(x);

    console.log(`encrypting the following value: ${value}`);

    return cipher.update(value, 'utf-8', 'hex') + cipher.final('hex');
  }

  decrypt(encrypted: string): T {
    console.log(`attempting to decrypt '${encrypted}'`);

    let decipher = crypto.createDecipheriv(
      encryptionAlgo,
      this.secretKey,
      this.initializationVectors.toString('hex').slice(0, 16),
    );
    let decrypted = decipher.update(encrypted, 'hex', 'utf-8');

    decrypted += decipher.final('utf8');

    console.log(`decrypted value '${decrypted}'`);

    return JSON.parse(decrypted);
  }
}
