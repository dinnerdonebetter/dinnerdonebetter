import { EncryptorDecryptor } from './index';

describe('basic', () => {
  it('should encrypt and decrypt', () => {
    const ed = new EncryptorDecryptor('HEREISA32BYTESECRETWHICHISMADEUP');

    const exampleInput = 'test';
    let encrypted = ed.encrypt(exampleInput);
    expect(encrypted).toHaveLength(32);

    let decrypted = ed.decrypt(encrypted);
    expect(decrypted).toEqual(exampleInput);
  });
});
