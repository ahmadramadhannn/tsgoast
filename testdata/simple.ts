// Simple TypeScript examples for testing

function greet(name: string): string {
    return `Hello, ${name}!`;
}

const add = (a: number, b: number): number => {
    return a + b;
};

interface User {
    id: number;
    name: string;
    email?: string;
}

type Point = {
    x: number;
    y: number;
};

export const PI = 3.14159;
