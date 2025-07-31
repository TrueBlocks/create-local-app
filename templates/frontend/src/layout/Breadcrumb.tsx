import { useActiveProject } from '@hooks';
import { Anchor, Breadcrumbs, Text } from '@mantine/core';
import { useLocation } from 'wouter';

export const Breadcrumb = () => {
  const [currentLocation, navigate] = useLocation();
  const isWizard = currentLocation.startsWith('/wizard');
  const { lastTab } = useActiveProject();
  const currentTab = lastTab[currentLocation];

  const pathnames = currentLocation.split('/').filter((x: string) => x);
  const breadcrumbItems = [
    { title: 'Home', path: '/' },
    ...pathnames.map((value: string, index: number) => {
      const path = `/${pathnames.slice(0, index + 1).join('/')}`;
      return { title: value.charAt(0).toUpperCase() + value.slice(1), path };
    }),
    ...(currentTab ? [{ title: currentTab, path: currentLocation }] : []),
  ];

  return (
    <Breadcrumbs
      separator=">"
      px="xl"
      style={{ marginTop: '.5rem', marginBottom: '.5rem' }}
    >
      {breadcrumbItems.map((item, index) => {
        const isHome = index === 0;
        const disabled = isHome && isWizard;
        if (disabled) {
          return (
            <Text key={0} size="md">
              Home
            </Text>
          );
        }

        return (
          <Anchor
            key={index}
            component="button"
            onClick={() => {
              if (!disabled) navigate(item.path);
            }}
            style={{
              cursor: disabled ? 'default' : 'pointer',
              opacity: disabled ? 0.4 : 1,
              pointerEvents: disabled ? 'none' : 'auto',
            }}
          >
            <Text size="md">{item.title}</Text>
          </Anchor>
        );
      })}
    </Breadcrumbs>
  );
};
