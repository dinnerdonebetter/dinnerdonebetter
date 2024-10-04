import { EncryptorDecryptor } from './index';

describe('basic', () => {
  it('should encrypt and decrypt strings', () => {
    const ed = new EncryptorDecryptor('HEREISA32BYTESECRETWHICHISMADEUP');

    const exampleInput = 'test';
    let encrypted = ed.encrypt(exampleInput);
    expect(encrypted).not.toBe('');

    let decrypted = ed.decrypt(encrypted);
    expect(decrypted).toEqual(exampleInput);
  });

  it('should encrypt and decrypt objects', () => {
    const ed = new EncryptorDecryptor('HEREISA32BYTESECRETWHICHISMADEUP');

    const expected = { things: 'stuff' };
    let encrypted = ed.encrypt(expected);
    expect(encrypted).not.toBe('');

    let decrypted = ed.decrypt(encrypted);
    expect(decrypted).toEqual(expected);
  });
});
