import React, { useEffect } from "react";
import * as PIXI from "pixi.js";
import { useState, useRef } from "react";

function Game() {
    const canvasRef = useRef(null);
    const appRef = useRef(null);
    const blockToDropRef = useRef(null);
    
    useEffect(() => {
        const initPixi = async() => {
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

            const oneBlock = new PIXI.Sprite(assets['1'].textures['default']);
            oneBlock.width = 100;
            oneBlock.height = 100;
            blockToDropRef.current = oneBlock;
            


            const oneBlockBreak = new PIXI.AnimatedSprite([
                assets['1'].textures['frame1'],
                assets['1'].textures['frame2'],
                assets['1'].textures['frame3'],
                assets['1'].textures['frame4']
            ]);

            app.stage.addChild(oneBlock);

            // oneBlockBreak.width = 100;
            // oneBlockBreak.height = 100;
            // oneBlockBreak.x = 300;
            // oneBlockBreak.y = 300;
            // oneBlockBreak.animationSpeed = 0.2;
            // oneBlockBreak.play();
            // app.stage.addChild(oneBlock);
            // app.stage.addChild(oneBlockBreak);

            
        
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

    useEffect(() => {
        const handleKeyDown = (e) => {
            switch (e.key) {
                case 'ArrowLeft': blockToDropRef.current.x -= 100; break;
                case 'ArrowRight': blockToDropRef.current.x += 100; break;
                default: break;
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