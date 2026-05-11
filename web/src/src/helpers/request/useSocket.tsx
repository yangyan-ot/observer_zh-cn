import { useCallback, useEffect, useRef, useState } from 'react';

import { connectSocket, ISocket } from './connectSocket';

export const useSocket = (options: ISocket, reconnect: boolean) => {
    const optionsRef = useRef<ISocket | null>(options);
    const socketRef = useRef<WebSocket | null>(null);
    const reconnectTimeoutRef = useRef<NodeJS.Timeout | null>(null);
    const mountedRef = useRef(true);

    useEffect(() => {
        optionsRef.current = options;
    }, [options]);

    const [readyState, setReadyState] = useState<WebSocket['readyState']>(WebSocket.CONNECTING);

    const sendMessage = useCallback((message: unknown, json?: boolean) => {
        if (socketRef.current && socketRef.current.readyState === WebSocket.OPEN) {
            socketRef.current.send(json ? JSON.stringify(message) : String(message));
        }
    }, []);

    useEffect(() => {
        mountedRef.current = true;

        const doConnect = () => {
            if (!optionsRef.current || socketRef.current || !mountedRef.current) {
                return;
            }

            socketRef.current = connectSocket({
                ...optionsRef.current,
                onClose: (e) => {
                    optionsRef.current?.onClose?.(e);
                    setReadyState(WebSocket.CLOSED);
                    if (reconnect && mountedRef.current) {
                        if (reconnectTimeoutRef.current) {
                            clearTimeout(reconnectTimeoutRef.current);
                        }
                        reconnectTimeoutRef.current = setTimeout(() => {
                            socketRef.current = null;
                            doConnect();
                        }, 1000);
                    } else {
                        socketRef.current = null;
                    }
                },
                onError: (e) => {
                    optionsRef.current?.onError?.(e);
                    setReadyState(WebSocket.CLOSED);
                    if (reconnect && mountedRef.current) {
                        if (reconnectTimeoutRef.current) {
                            clearTimeout(reconnectTimeoutRef.current);
                        }
                        reconnectTimeoutRef.current = setTimeout(() => {
                            socketRef.current = null;
                            doConnect();
                        }, 1000);
                    } else {
                        socketRef.current = null;
                    }
                },
                onOpen: () => {
                    setReadyState(WebSocket.OPEN);
                },
                onData: (data) => {
                    optionsRef.current?.onData?.(data);
                }
            });
        };

        doConnect();

        return () => {
            mountedRef.current = false;
            if (reconnectTimeoutRef.current) {
                clearTimeout(reconnectTimeoutRef.current);
            }
            socketRef.current?.close();
            socketRef.current = null;
        };
    }, [reconnect]);

    return { sendMessage, readyState };
};
