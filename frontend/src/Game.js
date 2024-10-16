import React, { useEffect } from "react";
import * as PIXI from "pixi.js";
import { useState, useRef } from "react";

function Game() {
    const canvasRef = useRef(null);
    const appRef = useRef(null);
    
    useEffect(() => {
        const initPixi = async() => {
            if (!canvasRef.current) return;

            const app = new PIXI.Application();
            await app.init({
                width: 700,
                height: 800,
                backgroundColor: 0xffffff
            });
            canvasRef.current.appendChild(app.canvas);
            appRef.current = app;

            /*
            PIXI.Assets.add([
                { alias: '1', src: '/sheets/1_block_sheet.json' },
                { alias: '2', src: '/sheets/2_block_sheet.json' },
                { alias: '3', src: '/sheets/3_block_sheet.json' },
                { alias: '4', src: '/sheets/4_block_sheet.json' },
                { alias: '5', src: '/sheets/5_block_sheet.json' },
                { alias: '6', src: '/sheets/6_block_sheet.json' },
                { alias: '7', src: '/sheets/7_block_sheet.json' },
                { alias: 'barrier', src: '/sheets/barrier_block_sheet.json' }
            ]);
            */
            //await PIXI.Assets.load(['1', '2', '3', '4', '5', '6', '7', 'barrier']);

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
        };
        initPixi();


        return () => {
            appRef.current.destroy(true, { children: true, texture: true, baseTexture: true });
            appRef.current = null;
        }
    });

    return (
        <canvas ref={canvasRef} />
    );
    
};

export default Game;