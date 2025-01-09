import * as fs from 'fs';
import * as path from 'path';

export function loadFile<T>(filepath: string): T {
  try {
    // Resolve the full path to ensure it works across environments
    const fullPath = path.resolve(filepath);

    // Read the file content
    const fileContent = fs.readFileSync(fullPath, 'utf-8');

    // Parse the JSON content
    return JSON.parse(fileContent) as T;
  } catch (error) {
    throw new Error(`Failed to load JSON from ${filepath}: ${error}`);
  }
}
