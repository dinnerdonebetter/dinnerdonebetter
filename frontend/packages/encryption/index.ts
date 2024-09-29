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

export class EncryptorDecryptor {
  secretKey: string;
  initializationVectors: Buffer;

  constructor(secretKey: string, initializationVectors?: Buffer) {
    if (!secretKey) {
      throw new Error('secretKey is required');
    }

    this.secretKey = secretKey;

    if (initializationVectors) {
      this.initializationVectors = initializationVectors;
    } else {
      this.initializationVectors = crypto.randomBytes(16);
    }
  }

  encrypt(x: string): string {
    let cipher = crypto.createCipheriv('aes-256-cbc', this.secretKey, this.initializationVectors);
    let encrypted = cipher.update(x, 'utf-8', 'hex');
    encrypted += cipher.final('hex');

    return encrypted
  }

  decrypt(encrypted: string): string {
    let decipher = crypto.createDecipheriv('aes-256-cbc', this.secretKey, this.initializationVectors);
    let decrypted = decipher.update(
     encrypted,
     'hex',
     'utf-8'
    );
    
    decrypted += decipher.final('utf8');
    return decrypted;
  }
}
