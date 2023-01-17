import { Typography } from '@mui/material';
import Grid from '@mui/material/Grid';
import Paper from '@mui/material/Paper';
import { ResponsiveBar } from '@nivo/bar';
import { ResponsivePie } from '@nivo/pie';
import DataTable from 'component/DataTable';
import { AuthContext } from 'context/auth';
import { useToast } from 'context/toast';
import moment from 'moment';
import { useContext, useEffect, useState } from 'react';
import { isAPIError } from 'service/error/model';
import { getAudits } from 'service/model/audit';
import { Stats } from 'service/model/model';
import { getStats } from 'service/model/stats';
import { AuditColumns } from 'page/dashboard/component/columns';
import { FilterRequest } from 'service/common/filter';
import { SortRequest } from 'service/common/sort';

type BarData = BarDataItem[]
type PieData = PieDataItem[]

type BarDataItem = {
  day: string
  allowed: number
  allowedColor: string
  denied: number
  deniedColor: string
}

type PieDataItem = {
  id: string
  value: number
  color: string
}

export default function Dashboard() {
  const { user } = useContext(AuthContext);
  const toast = useToast();

  const [stats, setStats] = useState<Stats>([]);

  const [allowedVsDeniedPerDayData, setAllowedVsDeniedPerDayData] = useState<BarData>([]);
  const [allowedVsDeniedData, setAllowedVsDeniedData] = useState<PieData>([]);

  useEffect(() => {
    if (user === undefined) {
      return;
    }

    const fetchStats = async () => {
      const response = await getStats(user?.token!);

      if (isAPIError(response)) {
        toast.error(`Unable to retrieve statistics: ${response.message}`);
      } else {
        setStats(response);
      }
    };

    fetchStats();
  // eslint-disable-next-line
  }, [user]);

  // Compute per day allowed vs denied check decisions data
  useEffect(() => {
    setAllowedVsDeniedPerDayData(
      stats.map(dayStat => {
        return {
          day: moment(dayStat.date).format('D MMM'),
          allowed: dayStat.checks_allowed_number,
          allowedColor: 'hsl(244, 70%, 50%)',
          denied: dayStat.checks_denied_number,
          deniedColor: 'hsl(130, 70%, 50%)',
        }
      }),
    );
  // eslint-disable-next-line
}, [stats]);

  // Compute allowed vs denied pie chart data
  useEffect(() => {
    const totalAllowed = stats.reduce((acc, cur) => acc + cur.checks_allowed_number, 0);
    const totalDenied = stats.reduce((acc, cur) => acc + cur.checks_denied_number, 0);

    setAllowedVsDeniedData([
      { id: 'Allowed', value: totalAllowed, color: 'hsl(244, 70%, 50%)' },
      { id: 'Denied', value: totalDenied, color: 'hsl(130, 70%, 50%)' },
    ]);
  }, [stats]);

  const auditColumns = AuditColumns();

  const fetcher = (page?: number, size?: number, filter?: FilterRequest, sort?: SortRequest) => {
    return getAudits(user?.token!, page, size, filter, sort);
  };

  return (
    <>
      <Grid container spacing={3}>
        <Grid item xs={12} md={6} lg={6}>
          <Paper sx={{ width: '100%', height: 340 }}>
            <Typography variant='h5' sx={{ pt: 2, pl: 2, pb: 0 }}>
              Checks decisions per day
            </Typography>
            <ResponsiveBar
                data={allowedVsDeniedPerDayData}
                keys={[
                    'allowed',
                    'denied',
                ]}
                colors={{ scheme: 'paired' }}
                indexBy='day'
                margin={{ top: 20, right: 20, bottom: 80, left: 60 }}
                padding={0.3}
                axisLeft={{
                    tickSize: 5,
                    tickPadding: 5,
                    tickRotation: 0,
                    legend: 'Allowed vs Denied',
                    legendPosition: 'middle',
                    legendOffset: -40
                }}
                
                legends={[
                  {
                      dataFrom: 'keys',
                      anchor: 'top',
                      direction: 'row',
                      justify: false,
                      translateX: 0,
                      translateY: -25,
                      itemsSpacing: 2,
                      itemWidth: 100,
                      itemHeight: 20,
                      itemDirection: 'left-to-right',
                      itemOpacity: 0.85,
                      symbolSize: 20,
                      effects: [
                          {
                              on: 'hover',
                              style: {
                                  itemOpacity: 1
                              }
                          }
                      ]
                  }
              ]}
            />
          </Paper>
        </Grid>

        <Grid item xs={12} md={6} lg={6}>
          <Paper sx={{ width: '100%', height: 340 }}>
            <Typography variant='h5' sx={{ pt: 2, pl: 2, pb: 0 }}>
              Total of check decisions
            </Typography>
            <ResponsivePie
                data={allowedVsDeniedData}
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

      <Grid container spacing={3} sx={{ mt: 0 }}>
        <DataTable
          title='Audit logs'
          defaultSize={5}
          columns={auditColumns}
          fetcher={fetcher}
          sx={{ p: 2 }}
        />
      </Grid>
    </>
  );
}