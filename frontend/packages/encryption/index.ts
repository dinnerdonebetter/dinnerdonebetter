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

export class EncryptorDecryptor<T> {
  secretKey: string;
  initializationVectors: Buffer;
  // debug, delete later
  initializedSecretKey: boolean = false;
  initializedIV: boolean = false;

  constructor(encryptionKey?: string, initializationVectors?: string) {
    if ((encryptionKey || '').length !== 32) {
      console.log('encryption key not provided, generating random bytes');
      this.secretKey = crypto.randomBytes(32).toString('hex').slice(0, 32);
    } else {
      console.log(`setting up encryption with provided encryptionKey ${encryptionKey}`);
      this.initializedSecretKey = true;
      this.secretKey = encryptionKey!;
    }

    if (Buffer.from(initializationVectors || '', 'base64').length !== 32) {
      console.log('initialization vectors not provided, generating random bytes');
      this.initializationVectors = crypto.randomBytes(32);
    } else {
      console.log(`setting up encryption with provided initializationVectors ${initializationVectors}`);
      this.initializedIV = true;
      this.initializationVectors = Buffer.from(initializationVectors!, 'base64');
    }
  }

  encrypt(x: T): string {
    console.log(
      `encrypting with initialization vectors: ${this.initializationVectors.toString('base64')}, initializedSecretKey: ${this.initializedSecretKey}, initializedIV: ${this.initializedIV}`,
    ); // TODO: DELETEME upon verification

    let cipher = crypto.createCipheriv(
      'aes-256-cbc',
      this.secretKey,
      this.initializationVectors.toString('hex').slice(0, 16),
    );
    let encrypted = cipher.update(JSON.stringify(x), 'utf-8', 'hex');
    encrypted += cipher.final('hex');

    return encrypted;
  }

  decrypt(encrypted: string): T {
    console.log(
      `decrypting with initialization vectors: ${this.initializationVectors.toString('base64')}, initializedSecretKey: ${this.initializedSecretKey}, initializedIV: ${this.initializedIV}`,
    ); // TODO: DELETEME upon verification

    let decipher = crypto.createDecipheriv(
      'aes-256-cbc',
      this.secretKey,
      this.initializationVectors.toString('hex').slice(0, 16),
    );
    let decrypted = decipher.update(encrypted, 'hex', 'utf-8');

    decrypted += decipher.final('utf8');

    return JSON.parse(decrypted);
  }
}
