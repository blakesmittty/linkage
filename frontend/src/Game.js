import React, { useEffect } from "react";
import * as PIXI from "pixi.js";
import { useState, useRef } from "react";


function Game() {
    const canvasRef = useRef(null);
    const appRef = useRef(null);
    const blockToDropRef = useRef(null);
    const blocksRef = useRef(null);
    const blockToDropNumRef = useRef(null);
    const [loaded, setLoaded] = useState(false);
    const [socket, setSocket] = useState(null);
    const socketRef = useRef(null);
    const oneRef = useRef(null);
    const twoRef = useRef(null);
    const threeRef = useRef(null);
    const fourRef = useRef()
    const occupiedPositionsRef = useRef(new Set()); 
    const spritesRef = useRef(new Set());
    
    useEffect(() => {
        const username = localStorage.getItem('username');
        const ws = new WebSocket(`ws://localhost:8080/ws?username=${username}&roomID=${generateRandomRoomID()}`);

        ws.onopen = () => {
            console.log("connection to server established");
        };

        ws.onerror = (error) => {
            console.log("error in websocket: ",error);
        };

        ws.onclose = () => {
            console.log("closing websocket connection");
        };

        ws.onmessage = (message) => {
            console.log("received game state: ", message.data);
            const data = JSON.parse(message.data);
            
            clearGrid();
            updateGrid(data.Grid);
        };

        setSocket(ws);
        socketRef.current = ws;

        const clearGrid = () => {
            for (const sprite of spritesRef.current) {
                appRef.current.stage.removeChild(sprite);
            }
            spritesRef.current.clear();
        };

        const occupySet = (grid) => {
            grid.forEach((row, rowIndex) => {
                row.forEach((cell, colIndex) => {
                    if (cell.Number != 0) {
                        const position = `${rowIndex},${colIndex}`;
                        occupiedPositionsRef.current.add(position);
                    }
                });
            });
        };

        const updateGrid = (grid) => {
            const app = appRef.current;
    
            if (!app || !app.stage) {
                console.log("PIXI application is not yet initialized");
                return;
            }
    
            const blocks = blocksRef.current;
    
            grid.forEach((row, rowIndex) => {
                row.forEach((cell, colIndex) => {
                    if (cell.Number == -2 ) {
                        const blockAsset = blocks[8]['default'];
                        const position = `${rowIndex},${colIndex}`;
                        const block = new PIXI.Sprite(blockAsset);
                        spritesRef.current.add(block);
                        block.x = colIndex * 100;
                        block.y = rowIndex * 100;
                        block.width = 100;
                        block.height = 100;
                        app.stage.addChild(block);
                    } else if (cell.Number == -1) {
                        const blockAsset = blocks[8]['broken'];
                        const position = `${rowIndex},${colIndex}`;
                        const block = new PIXI.Sprite(blockAsset);
                        spritesRef.current.add(block);
                        block.x = colIndex * 100;
                        block.y = rowIndex * 100;
                        block.width = 100;
                        block.height = 100;
                        app.stage.addChild(block);

                    }
                    if (cell.Number > 0) {
                        const blockAsset = blocks[cell.Number]['default'];
                        const position = `${rowIndex},${colIndex}`;
                        if (!occupiedPositionsRef.current.has(position)) {
                            const block = new PIXI.Sprite(blockAsset);
                            spritesRef.current.add(block);
                            block.x = colIndex * 100;
                            block.y = rowIndex * 100;
                            block.width = 100;
                            block.height = 100;
                            app.stage.addChild(block);
                        }
                    }
                });
            });
        };

        return () => {
            if (ws.readyState === WebSocket.OPEN) {
              ws.close();
            }
          };

    }, []);

    useEffect(() => {

        const initPixi = async () => {
            if (!canvasRef.current) return;

            const app = new PIXI.Application();
            await app.init({
                width: 700,
                height: 800,
                backgroundColor: 'black'
            });
            canvasRef.current.appendChild(app.canvas);
            appRef.current = app;
            console.log("loading assets");

            PIXI.Assets.addBundle('blocks', [
                { alias: '1', src: '/sheets/1_block_sheet.json' },
                { alias: '2', src: '/sheets/2_block_sheet.json' },
                { alias: '3', src: '/sheets/3_block_sheet.json' },
                { alias: '4', src: '/sheets/4_block_sheet.json' },
                { alias: '5', src: '/sheets/5_block_sheet.json' },
                { alias: '6', src: '/sheets/6_block_sheet.json' },
                { alias: '7', src: '/sheets/7_block_sheet.json' },
                { alias: 'barrier', src: '/sheets/barrier_block_sheet.json' }
            ]);

            const assets = await PIXI.Assets.loadBundle('blocks');
            console.log('assets loaded');

            const block1= assets['1'].textures['default'];
            const block2 = assets['2'].textures['default'];
            const block3 = assets['3'].textures['default'];
            const block4 = assets['4'].textures['default'];
            const block5 = assets['5'].textures['default'];
            const block6 = assets['6'].textures['default'];
            const block7 = assets['7'].textures['default'];
            const blockBarrier = assets['barrier'].textures['default'];
            const blockBarrierBroken = assets['barrier'].textures['frame1'];

            // block1.width = 100;
            // block1.height = 100;
            // blockToDropRef.current = block1;

            const block1Break= new PIXI.AnimatedSprite([
                assets['1'].textures['frame1'],
                assets['1'].textures['frame2'],
                assets['1'].textures['frame3'],
                assets['1'].textures['frame4']
            ]);

            const block2Break = new PIXI.AnimatedSprite([
                assets['2'].textures['frame1'],
                assets['2'].textures['frame2'],
                assets['2'].textures['frame3'],
                assets['2'].textures['frame4']
            ]);

            const block3Break = new PIXI.AnimatedSprite([
                assets['3'].textures['frame1'],
                assets['3'].textures['frame2'],
                assets['3'].textures['frame3'],
                assets['3'].textures['frame4']
            ]);

            const block4Break = new PIXI.AnimatedSprite([
                assets['4'].textures['frame1'],
                assets['4'].textures['frame2'],
                assets['4'].textures['frame3'],
                assets['4'].textures['frame4']
            ]);

            const block5Break = new PIXI.AnimatedSprite([
                assets['5'].textures['frame1'],
                assets['5'].textures['frame2'],
                assets['5'].textures['frame3'],
                assets['5'].textures['frame4']
            ]);

            const block6Break = new PIXI.AnimatedSprite([
                assets['6'].textures['frame1'],
                assets['6'].textures['frame2'],
                assets['6'].textures['frame3'],
                assets['6'].textures['frame4']
            ]);

            const block7Break = new PIXI.AnimatedSprite([
                assets['7'].textures['frame1'],
                assets['7'].textures['frame2'],
                assets['7'].textures['frame3'],
                assets['7'].textures['frame4']
            ]);

            const blocks = {
                1: {
                    'default': block1,
                    'break': block1Break
                },
                2: {
                    'default': block2,
                    'break': block2Break
                },
                3: {
                    'default': block3,
                    'break': block3Break
                },
                4: {
                    'default': block4,
                    'break': block4Break
                },
                5: {
                    'default': block5,
                    'break': block5Break
                },
                6: {
                    'default': block6,
                    'break': block6Break
                },
                7: {
                    'default': block7,
                    'break': block7Break
                },
                8: {
                    'default': blockBarrier,
                    'broken': blockBarrierBroken
                }   
            };
            blocksRef.current = blocks;

            let blockNum = getRandomInt(1, 7);
            let randBlockAsset = blocks[blockNum]['default'];
            const randBlock = new PIXI.Sprite(randBlockAsset);
            //spritesRef.current.add(randBlock);
            console.log(randBlock);
            blockToDropNumRef.current = blockNum;
            blockToDropRef.current = randBlock;

            
            randBlock.x = 300;
            randBlock.y = 0;
            randBlock.width = 100;
            randBlock.height = 100;

            app.stage.addChild(randBlock);

            // oneBlockBreak.width = 100;
            // oneBlockBreak.height = 100;
            // oneBlockBreak.x = 300;
            // oneBlockBreak.y = 300;
            // oneBlockBreak.animationSpeed = 0.2;
            // oneBlockBreak.play();
            // app.stage.addChild(oneBlock);
            // app.stage.addChild(oneBlockBreak);

            
            setLoaded(true);
        };
        initPixi();


        return () => {
            if (appRef.current) {
                console.log("destroying pixi app");
                appRef.current.destroy(true, { children: true, texture: true, baseTexture: true });
                appRef.current = null;
            }
        }
    }, []);


    function generateRandomRoomID(length = 8) {
        const characters = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789';
        let roomID = '';
        for (let i = 0; i < length; i++) {
            const randomIndex = Math.floor(Math.random() * characters.length);
            roomID += characters[randomIndex];
        }
        return roomID;
    };
   
    function getRandomInt(min, max) {
        min = Math.ceil(min);
        max = Math.floor(max);
        return Math.floor(Math.random() * (max - min + 1)) + min;
    }

    useEffect(() => {
        const sendPlayerInput = () => {
            const block = blockToDropRef.current;
            const column = Math.floor(block.x / 100);
            
            const socket = socketRef.current;
            if (socket && socket.readyState === WebSocket.OPEN) {
                const message = JSON.stringify({
                    action: "drop",
                    block: blockToDropNumRef.current,
                    column: column
                });
                socket.send(message);
            } else {
                console.log("socket not ready");
            }
            appRef.current.stage.removeChild(blockToDropRef.current);

            let blockNum = getRandomInt(1, 7);
            let randBlockAsset = blocksRef.current[blockNum]['default'];
            const randBlock = new PIXI.Sprite(randBlockAsset);
           // spritesRef.current.add(randBlock);

            randBlock.x = 300;
            randBlock.y = 0;
            randBlock.width = 100;
            randBlock.height = 100
            appRef.current.stage.addChild(randBlock);
            blockToDropRef.current = randBlock;
            blockToDropNumRef.current = blockNum;
        };

        const handleKeyDown = (e) => {
            const block = blockToDropRef.current;
            const maxX = 700 - block.width;
            const step = 100;
            const minX = 0;
            const socket = socketRef.current;
            if (socket && socket.readyState === WebSocket.OPEN) {
                switch (e.key) {
                    case 'a':
                    case 'ArrowLeft': block.x = Math.max(block.x - step, minX); break;
                    case 'd':
                    case 'ArrowRight': block.x = Math.min(block.x + step, maxX); break;
                    case 's': 
                    case 'ArrowDown': sendPlayerInput(); break;
                    default: break;
                }
            } else {
                console.log("socket not ready");
            }
        }

        window.addEventListener('keydown', handleKeyDown);

        return () => {
            window.removeEventListener('keydown', handleKeyDown);
        }
    }, []);

    return (
        <div ref={canvasRef} />
    );
    
};

export default Game;