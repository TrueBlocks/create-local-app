import { useEffect } from 'react';

import { NodeStatus, getBarWidth } from '@components';
import { ViewContextProvider, WalletConnectProvider } from '@contexts';
import {
  useActiveProject,
  useAppHealth,
  useAppHotkeys,
  useAppNavigation,
} from '@hooks';
import { Footer, Header, HelpBar, MainView, MenuBar } from '@layout';
import { AppShell } from '@mantine/core';
import { WalletConnectModalSign } from '@walletconnect/modal-sign-react';
import { Router } from 'wouter';

import { useGlobalEscape } from './hooks/useGlobalEscape';

// Add at the top level, outside the component
function globalNavKeySquelcher(e: KeyboardEvent) {
  const navKeys = [
    'ArrowUp',
    'ArrowDown',
    'ArrowLeft',
    'ArrowRight',
    'PageUp',
    'PageDown',
    'Home',
    'End',
  ];

  const activeElement = document.activeElement as HTMLElement;
  const isFormElement =
    activeElement &&
    (activeElement.tagName === 'INPUT' ||
      activeElement.tagName === 'TEXTAREA' ||
      activeElement.tagName === 'SELECT' ||
      activeElement.isContentEditable);

  if (navKeys.includes(e.key) && !isFormElement) {
    // Only squelch if not handled by a focused form control
    e.preventDefault();
  }
}

export const App = () => {
  useEffect(() => {
    window.addEventListener('keydown', globalNavKeySquelcher, {
      capture: true,
    });
    return () => {
      window.removeEventListener('keydown', globalNavKeySquelcher, {
        capture: true,
      });
    };
  }, []);

  const { ready, isWizard } = useAppNavigation();
  const { menuCollapsed, helpCollapsed } = useActiveProject();

  useAppHotkeys();
  useAppHealth();
  useGlobalEscape();

  if (!ready) return <div>Not ready</div>;

  const header = { height: 60 };
  const footer = { height: 40 };
  const navbar = {
    width: getBarWidth(menuCollapsed, 1),
    breakpoint: 'sm',
    collapsed: { mobile: !menuCollapsed },
  };
  const aside = {
    width: getBarWidth(helpCollapsed, 2),
    breakpoint: 'sm',
    collapsed: { mobile: !helpCollapsed },
  };

  return (
    <Router>
      <WalletConnectProvider>
        <div
          style={{
            display: 'flex',
            flexDirection: 'column',
            height: '100vh',
          }}
        >
          <AppShell
            layout="default"
            header={header}
            footer={footer}
            navbar={navbar}
            aside={aside}
          >
            <Header />
            <MenuBar disabled={isWizard} />
            <ViewContextProvider>
              <MainView />
            </ViewContextProvider>
            <HelpBar />
            <div
              style={{
                position: 'absolute',
                top: '84px',
                right: `${getBarWidth(helpCollapsed, 2) + 0}px`,
                zIndex: 1000,
              }}
            >
              <NodeStatus />
            </div>
            <Footer />
          </AppShell>
          <WalletConnectModalSign
            projectId={
              import.meta.env.VITE_WALLETCONNECT_PROJECT_ID ||
              (() => {
                console.error(
                  'VITE_WALLETCONNECT_PROJECT_ID not set in environment variables',
                );
                return 'MISSING_PROJECT_ID';
              })()
            }
            metadata={{
              name: '{{ORG_NAME}} {{PROJECT_PROPER}}',
              description:
                'A {{ORG_NAME}} desktop application for naming addresses',
              url: 'https://{{DOMAIN}}',
              icons: ['https://{{DOMAIN}}/favicon.ico'],
            }}
          />
        </div>
      </WalletConnectProvider>
    </Router>
  );
};
