import { EncryptorDecryptor } from './index';

describe('basic', () => {
  it('should encrypt and decrypt strings', () => {
    const ed = new EncryptorDecryptor(
      'HEREISA32BYTESECRETWHICHISMADEUP',
      Buffer.from('HEREISA32BYTESECRETWHICHISMADEUP').toString('base64'),
    );

    const exampleInput = 'test';
    let encrypted = ed.encrypt(exampleInput);
    expect(encrypted).not.toBe('');

    let decrypted = ed.decrypt(encrypted);
    expect(decrypted).toEqual(exampleInput);
  });

  it('should encrypt and decrypt strings without params', () => {
    const ed = new EncryptorDecryptor();

    const exampleInput = 'test';
    let encrypted = ed.encrypt(exampleInput);
    expect(encrypted).not.toBe('');

    let decrypted = ed.decrypt(encrypted);
    expect(decrypted).toEqual(exampleInput);
  });

  it('should encrypt and decrypt objects', () => {
    const ed = new EncryptorDecryptor(
      'HEREISA32BYTESECRETWHICHISMADEUP',
      Buffer.from('HEREISA32BYTESECRETWHICHISMADEUP').toString('base64'),
    );

    const expected = { things: 'stuff' };
    let encrypted = ed.encrypt(expected);
    expect(encrypted).not.toBe('');

    let decrypted = ed.decrypt(encrypted);
    expect(decrypted).toEqual(expected);
  });

  it('should encrypt and decrypt objects without params', () => {
    const ed = new EncryptorDecryptor();

    const expected = { things: 'stuff' };
    let encrypted = ed.encrypt(expected);
    expect(encrypted).not.toBe('');

    let decrypted = ed.decrypt(encrypted);
    expect(decrypted).toEqual(expected);
  });
});
