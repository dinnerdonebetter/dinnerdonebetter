import * as crypto from 'crypto';

export class ServerTimingEvent {
  name: string = '';
  description: string = '';
  startTime: Date = new Date();
  endTime: Date = new Date();
  duration: number = 0;

  constructor(name: string = '', description: string = '') {
    this.name = name;
    this.description = description;
    this.startTime = new Date();
  }

  end() {
    this.endTime = new Date();
    this.duration = this.endTime.getTime() - this.startTime.getTime();
  }
}

export class EncryptorDecryptor<JSONValue> {
  secretKey: string;
  initializationVectors: Buffer;

  constructor(secretKey: string, initializationVectors?: string) {
    if (!secretKey) {
      throw new Error('secretKey is required');
    }

    this.secretKey = secretKey;

    if (initializationVectors) {
      this.initializationVectors = Buffer.from(initializationVectors, 'base64');
    } else {
      this.initializationVectors = crypto.randomBytes(16);
    }

    console.log(`initializationVectors: ${this.initializationVectors.toString('base64')}`);
  }

  encrypt(x: JSONValue): string {
    let cipher = crypto.createCipheriv('aes-256-cbc', this.secretKey, this.initializationVectors.toString('hex').slice(0, 16));
    let encrypted = cipher.update(JSON.stringify(x), 'utf-8', 'hex');
    encrypted += cipher.final('hex');

    return encrypted;
  }

  decrypt(encrypted: string): JSONValue {
    let decipher = crypto.createDecipheriv('aes-256-cbc', this.secretKey, this.initializationVectors.toString('hex').slice(0, 16));
    let decrypted = decipher.update(encrypted, 'hex', 'utf-8');

    decrypted += decipher.final('utf8');

    return JSON.parse(decrypted);
  }
}
