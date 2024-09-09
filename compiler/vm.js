let PUSH = 'PUSH'
let ADD = 'ADD'
let MINUS = 'MINUS'

let vm = function (program) {
    let programCounter = 0
    let stackPointer = 0;
    let stack = []

    while (programCounter < program.length) {
        let instruction = program[programCounter]
        programCounter++
        let left, right
        switch (instruction) {
            case PUSH:
                stack[stackPointer] = program[programCounter]
                programCounter++
                stackPointer++
                break
            case ADD:
                right = stack[stackPointer - 1]
                stackPointer--
                left = stack[stackPointer - 1]
                stack[stackPointer - 1] = left + right
                break
            case MINUS:
                right = stack[stackPointer - 1]
                stackPointer--
                left = stack[stackPointer - 1]
                stack[stackPointer - 1] = left - right
                break
        }
    }

    console.log(stack[stackPointer - 1])
}

let program = [
    PUSH, 3,
    PUSH, 4,
    ADD,
    PUSH, 5,
    MINUS,
    PUSH, 5,
    PUSH, 10,
    ADD,
    ADD
];

vm(program)
