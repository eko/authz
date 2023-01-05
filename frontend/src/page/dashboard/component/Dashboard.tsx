import { Typography } from '@mui/material';
import Grid from '@mui/material/Grid';
import Paper from '@mui/material/Paper';
import { ResponsiveBar } from '@nivo/bar';
import { ResponsivePie } from '@nivo/pie';

const barData = [
  {
    "provider": "One",
    "success": 74,
    "successColor": "hsl(152, 100%, 42%)",
    "error": 111,
    "errorColor": "hsl(297, 70%, 50%)",
  },
  {
    "provider": "Google Cloud Platform",
    "success": 74,
    "successColor": "hsl(152, 100%, 42%)",
    "error": 111,
    "errorColor": "hsl(297, 70%, 50%)",
  },
  {
    "provider": "Two",
    "success": 74,
    "successColor": "hsl(152, 100%, 42%)",
    "error": 111,
    "errorColor": "hsl(297, 70%, 50%)",
  },
];

const pieData = [
  {
    "id": "allowed",
    "label": "allowed",
    "value": 139,
    "color": "hsl(244, 70%, 50%)"
  },
  {
    "id": "denied",
    "label": "denied",
    "value": 508,
    "color": "hsl(130, 70%, 50%)"
  }
];

export default function Dashboard() {
  return (
    <Grid container spacing={3}>
      <Grid item xs={12} md={6} lg={6}>
        <Paper sx={{ width: '100%', height: 340 }}>
          <Typography variant='h5' sx={{ pt: 2, pl: 2, pb: 0 }}>
            Test
          </Typography>
          <ResponsiveBar
              data={barData}
              keys={[
                  'success',
                  'error',
              ]}
              colors={{ scheme: 'paired' }}
              indexBy="provider"
              margin={{ top: 20, right: 20, bottom: 80, left: 60 }}
              padding={0.3}
              axisLeft={{
                  tickSize: 5,
                  tickPadding: 5,
                  tickRotation: 0,
                  legend: "Test",
                  legendPosition: 'middle',
                  legendOffset: -40
              }}
          />
        </Paper>
      </Grid>

      <Grid item xs={12} md={6} lg={6}>
        <Paper sx={{ width: '100%', height: 340 }}>
          <Typography variant='h5' sx={{ pt: 2, pl: 2, pb: 0 }}>
            Allowed vs Denied
          </Typography>
          <ResponsivePie
              data={pieData}
              colors={{ scheme: 'paired' }}
              margin={{ top: 20, right: 60, bottom: 80, left: 60 }}
              innerRadius={0.5}
              cornerRadius={3}
              activeOuterRadiusOffset={8}
              arcLinkLabelsSkipAngle={10}
              arcLinkLabelsTextColor="#333333"
              arcLinkLabelsThickness={2}
              arcLinkLabelsColor={{ from: 'color' }}
              arcLabelsSkipAngle={10}
          />
        </Paper>
      </Grid>
    </Grid>
  );
}