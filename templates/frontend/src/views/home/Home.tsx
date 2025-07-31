import { Suspense, useEffect, useState } from 'react';

import { GetAppId } from '@app';
import { LazyPanel, PanelSkeleton } from '@components';
import { Container, Grid, Space, Text, Title } from '@mantine/core';
import { Log } from '@utils';
import { useLocation } from 'wouter';

import styles from './Home.module.css';
import {
  LazyAbisPanel,
  LazyChunksPanel,
  LazyExportsPanel,
  LazyMonitorsPanel,
  LazyNamesPanel,
  LazyProjectsPanel,
} from './LazyPanels';
import { QuickActions } from './QuickActions';
import { RecentActivity } from './RecentActivity';
import { SampleImageSection } from './SampleImageSection';
import { useUploadAddressesDialog } from './UploadAddressesDialog';

export const Home = () => {
  const [_, setAppName] = useState('Your App');
  const [, navigate] = useLocation();
  const uploadDialog = useUploadAddressesDialog();

  useEffect(() => {
    GetAppId().then((id) => {
      setAppName(id.appName);
    });
  }, []);

  const handleAddressesUpload = async (addresses: string[]) => {
    // TODO: Implement actual address upload to monitors
    Log(`Uploading addresses not implemented:, ${addresses}`);
    // This would typically call an API to add multiple monitors
    // For now, we'll just log and navigate to monitors page
    navigate('/monitors');
  };

  const handleViewGallery = () => {
    // TODO: Implement gallery view or navigate to output folder
    Log('Navigate to gallery not implemented');
  };

  return (
    <Container size="xl" py="xl" className={styles.homeContainer}>
      <Space h="md" />

      {/* Header */}
      <header role="banner">
        <Title order={1} mb="xs">
          {{ORG_NAME}} DalleDress
        </Title>
        <Text size="lg" c="dimmed" mb="xl">
          Generate stunning images from Ethereum addresses, block hashes, or
          transaction hashes.
        </Text>
      </header>

      <main role="main" aria-label="Dashboard Overview">
        <Grid className={styles.dashboardGrid}>
          <Grid.Col span={{ base: 12, md: 8 }}>
            <section aria-label="Data Overview Panels" aria-live="polite">
              <Grid>
                <PanelProjects />
                <PanelExports />
                <PanelMonitors />
                <PanelChunks />
                <PanelNames />
                <PanelAbis />
              </Grid>
            </section>
          </Grid.Col>

          <Grid.Col span={{ base: 12, md: 4 }}>
            <aside
              role="complementary"
              aria-label="Quick Actions and Recent Activity"
            >
              <div className={styles.rightColumn}>
                <PanelRecentActivity />
                <PanelQuickActions viewGallery={handleViewGallery} />
                <PanelImageSection viewGallery={handleViewGallery} />
              </div>
            </aside>
          </Grid.Col>
        </Grid>
      </main>

      <uploadDialog.Dialog
        opened={uploadDialog.opened}
        onClose={uploadDialog.close}
        onUpload={handleAddressesUpload}
      />
    </Container>
  );
};

const PanelNames = () => {
  const [, navigate] = useLocation();
  return (
    <Grid.Col span={{ base: 12, sm: 6 }}>
      <div className={styles.panelCard}>
        <LazyPanel priority="high">
          <Suspense fallback={<PanelSkeleton title="Names" />}>
            <LazyNamesPanel
              onViewAll={() => {
                navigate('/names');
              }}
            />
          </Suspense>
        </LazyPanel>
      </div>
    </Grid.Col>
  );
};

const PanelMonitors = () => {
  const [, navigate] = useLocation();
  return (
    <Grid.Col span={{ base: 12, sm: 6 }}>
      <div className={styles.panelCard}>
        <LazyPanel priority="normal">
          <Suspense fallback={<PanelSkeleton title="Monitors" />}>
            <LazyMonitorsPanel
              onViewAll={() => {
                navigate('/monitors');
              }}
              onAddMonitor={() => Log('Add monitor not implemented')}
            />
          </Suspense>
        </LazyPanel>
      </div>
    </Grid.Col>
  );
};

const PanelExports = () => {
  const [, navigate] = useLocation();
  return (
    <Grid.Col span={{ base: 12, sm: 6 }}>
      <div className={styles.panelCard}>
        <LazyPanel priority="normal">
          <Suspense fallback={<PanelSkeleton title="Exports" />}>
            <LazyExportsPanel
              onViewAll={() => {
                navigate('/exports');
              }}
            />
          </Suspense>
        </LazyPanel>
      </div>
    </Grid.Col>
  );
};

const PanelProjects = () => {
  const [, navigate] = useLocation();
  return (
    <Grid.Col span={{ base: 12, sm: 6 }}>
      <div className={styles.panelCard}>
        <LazyPanel priority="high">
          <Suspense fallback={<PanelSkeleton title="Projects" />}>
            <LazyProjectsPanel
              onViewAll={() => {
                // Navigate to home page (projects are shown there)
                navigate('/');
              }}
            />
          </Suspense>
        </LazyPanel>
      </div>
    </Grid.Col>
  );
};

const PanelChunks = () => {
  const [, navigate] = useLocation();
  return (
    <Grid.Col span={{ base: 12, sm: 6 }}>
      <div className={styles.panelCard}>
        <LazyPanel priority="low">
          <Suspense fallback={<PanelSkeleton title="Chunks" />}>
            <LazyChunksPanel
              onViewAll={() => {
                navigate('/chunks');
              }}
            />
          </Suspense>
        </LazyPanel>
      </div>
    </Grid.Col>
  );
};

const PanelAbis = () => {
  const [, navigate] = useLocation();
  return (
    <Grid.Col span={{ base: 12, sm: 6 }}>
      <div className={styles.panelCard}>
        <LazyPanel priority="low">
          <Suspense fallback={<PanelSkeleton title="Abis" />}>
            <LazyAbisPanel
              onViewAll={() => {
                navigate('/abis');
              }}
            />
          </Suspense>
        </LazyPanel>
      </div>
    </Grid.Col>
  );
};

const PanelRecentActivity = () => {
  const [, navigate] = useLocation();
  const mockActivities = [
    {
      id: '1',
      type: 'image' as const,
      message: 'Generated new image from address 0x1234...5678',
      timestamp: new Date(Date.now() - 1000 * 60 * 30),
      onClick: () => navigate('/{{PROJECT_NAME}}'),
    },
    {
      id: '2',
      type: 'monitor' as const,
      message: 'Added monitor for 0xabcd...efgh',
      timestamp: new Date(Date.now() - 1000 * 60 * 60 * 2),
      onClick: () => navigate('/monitors'),
    },
    {
      id: '3',
      type: 'export' as const,
      message: 'Exported transaction data as CSV',
      timestamp: new Date(Date.now() - 1000 * 60 * 60 * 24),
      onClick: () => navigate('/exports'),
    },
  ];

  return (
    <section aria-label="Recent Activity" role="region">
      <div className={styles.activityContainer}>
        <LazyPanel priority="low">
          <RecentActivity activities={mockActivities} />
        </LazyPanel>
      </div>
    </section>
  );
};

type GalleryProps = {
  viewGallery: () => void;
};

const PanelQuickActions = ({ viewGallery }: GalleryProps) => {
  const [, navigate] = useLocation();
  const uploadDialog = useUploadAddressesDialog();

  const handleGenerateImage = () => {
    navigate('/{{PROJECT_NAME}}');
  };

  const handleUploadAddresses = () => {
    uploadDialog.open();
  };

  return (
    <section aria-label="Quick Actions" role="region">
      <div className={styles.quickActionsContainer}>
        <QuickActions
          onGenerateImage={handleGenerateImage}
          onViewGallery={viewGallery}
          onUploadAddresses={handleUploadAddresses}
        />
      </div>
    </section>
  );
};

const PanelImageSection = ({ viewGallery }: GalleryProps) => {
  return (
    <section aria-label="Sample Image" role="region">
      <div className={styles.sampleImageContainer}>
        <SampleImageSection onViewGallery={viewGallery} />
      </div>
    </section>
  );
};
