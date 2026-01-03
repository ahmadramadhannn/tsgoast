// Function test cases

// Regular function
function regularFunction(x: number, y: number): number {
    return x + y;
}

// Arrow function
const arrowFunc = (x: number): number => x * 2;

// Async function
async function fetchData(url: string): Promise<string> {
    const response = await fetch(url);
    return response.text();
}

// Async arrow function
const asyncArrow = async (id: number): Promise<void> => {
    console.log(id);
};

// Function with optional parameters
function withOptional(required: string, optional?: number): void {
    console.log(required, optional);
}

// Function with default parameters
function withDefault(name: string, age: number = 18): void {
    console.log(name, age);
}

// Function with rest parameters
function sum(...numbers: number[]): number {
    return numbers.reduce((a, b) => a + b, 0);
}

// Generic function
function identity<T>(arg: T): T {
    return arg;
}

// Exported function
export function exportedFunc(): void {
    console.log("exported");
}

// Generator function
function* generator(): Generator<number> {
    yield 1;
    yield 2;
    yield 3;
}
