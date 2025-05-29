'use server'
import { exec } from "node:child_process";

export async function runExe(filePath: string): Promise<string> {
    return new Promise((resolve, reject) => {
        exec(`start ${filePath}`, (error, stdout, stderr) => {
            if (error) {
                reject(error);
                return;
            }
            if (stderr) {
                reject(stderr)
                return
            }
            console.log("started controller file")
            resolve(stdout);
        });
    });
}