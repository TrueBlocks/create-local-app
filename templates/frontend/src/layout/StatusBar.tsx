import { useEffect, useState } from 'react';

import { useEvent } from '@hooks';
import { msgs } from '@models';

import './StatusBar.css';

export const StatusBar = () => {
  const [status, setStatus] = useState('');
  const [visible, setVisible] = useState(false);
  const [cn, setCn] = useState('okay');

  useEvent(msgs.EventType.STATUS, (message: string) => {
    setCn('okay');
    setStatus(message);
    setVisible(true);
  });

  useEvent(msgs.EventType.ERROR, (message: string) => {
    setCn('error');
    setStatus(message);
    setVisible(true);
  });

  useEffect(() => {
    if (!visible) return;
    const timer = setTimeout(() => {
      setVisible(false);
    }, 1500);
    return () => clearTimeout(timer);
  }, [visible, status]);

  if (!visible) return null;

  return (
    <div className={cn}>
      <span>{status}</span>
    </div>
  );
};
