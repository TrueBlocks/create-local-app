import { useEffect, useRef, useState } from 'react';

import { GetNodeStatus } from '@app';
import { types } from '@models';
import { useLocation } from 'wouter';

export const NodeStatus = () => {
  const [status, setStatus] = useState('Loading...');
  const [spinnerIndex, setSpinnerIndex] = useState(0);
  const [location, navigate] = useLocation();
  const previousView = useRef<string | null>(null);

  useEffect(() => {
    const getStatus = async () => {
      try {
        GetNodeStatus().then((meta: types.MetaData) => {
          const spinnerFrames = '⠋⠙⠹⠸⠼⠴⠦⠧⠇⠏';
          const dist = meta.client - meta.unripe;
          const formatter = new Intl.NumberFormat(navigator.language);
          const formattedClient = formatter.format(meta.client);
          const formattedUnripe = formatter.format(meta.unripe);
          const formattedDist = formatter.format(dist);
          const frame = spinnerFrames[spinnerIndex % spinnerFrames.length];
          setSpinnerIndex((prevIndex) => prevIndex + 1);
          const statusMessage = `${meta.chain}: ${frame} ${formattedClient} / ${formattedUnripe} / ${dist < 6 ? 'caught up' : formattedDist} `;
          setStatus(statusMessage);
        });
      } catch (error) {
        setStatus('Error fetching status ' + error);
      }
    };
    getStatus();
    const interval = setInterval(() => {
      getStatus();
    }, 1500);
    return () => {
      clearInterval(interval);
    };
  }, [spinnerIndex]);

  useEffect(() => {
    if (location !== '/khedra') {
      previousView.current = location;
    }
  }, [location]);

  const handleClick = () => {
    if (location === '/khedra' && previousView.current) {
      navigate(previousView.current);
    } else {
      navigate('/khedra');
    }
  };

  return (
    <div
      onClick={handleClick}
      style={{
        position: 'absolute',
        bottom: 0,
        right: 0,
        padding: '0px 10px',
        whiteSpace: 'nowrap',
        textAlign: 'right',
        zIndex: 9999,
        cursor: 'pointer',
        backgroundColor: 'rgba(100, 100, 220, 0.2)',
      }}
    >
      {status}
    </div>
  );
};
