import { useCallback, useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Rocket, Calendar, Newspaper, MapPin, Building2, X, Layers, Target } from 'lucide-react';
import { statsService, taskService } from '@/services';
import type { DashboardStats } from '@/services/statsService';
import { Button } from '@/components/ui/button';
import type { TaskAction, TaskInfo, TaskType } from '@/services/taskService';

const TASK_TYPE_OPTIONS: Array<{ type: TaskType; label: string }> = [
  { type: 'launch', label: 'Launch' },
  { type: 'agency', label: 'Agency' },
  { type: 'launcher', label: 'Launcher' },
  { type: 'launcher_family', label: 'Launcher Family' },
  { type: 'pad', label: 'Pad' },
  { type: 'location', label: 'Location' },
  { type: 'upcoming', label: 'Upcoming' },
];

export default function Dashboard() {
  const navigate = useNavigate();
  const [stats, setStats] = useState<DashboardStats>({
    totalRockets: 0,
    upcomingLaunches: 0,
    newsArticles: 0,
    launchBases: 0,
    companies: 0,
    ll2: {
      launches: 0,
      agencies: 0,
      launchers: 0,
      launcher_families: 0,
      locations: 0,
      pads: 0,
    },
  });
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [task, setTask] = useState<TaskInfo | null>(null);
  const [loadingTask, setLoadingTask] = useState(false);
  const [taskPendingAction, setTaskPendingAction] = useState<string | null>(null);
  const [showWelcome, setShowWelcome] = useState(
    () => localStorage.getItem('launchdate-admin-welcome-closed') !== 'true'
  );

  useEffect(() => {
    fetchStats();
  }, []);

  const fetchTask = useCallback(async () => {
    try {
      setLoadingTask(true);
      const currentTask = await taskService.getCurrentTask();
      setTask(currentTask);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to fetch task info');
      setTask(null);
    } finally {
      setLoadingTask(false);
    }
  }, []);

  useEffect(() => {
    fetchTask();
    const timer = window.setInterval(fetchTask, 10000);
    return () => window.clearInterval(timer);
  }, [fetchTask]);

  const fetchStats = async () => {
    try {
      setLoading(true);
      setError(null);
      const data = await statsService.getStats();
      setStats(data);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to fetch statistics');
    } finally {
      setLoading(false);
    }
  };

  const handleTaskAction = async (action: TaskAction) => {
    try {
      setTaskPendingAction(action);
      await taskService.actionTask(action);
      await fetchTask();
    } catch (err) {
      setError(err instanceof Error ? err.message : `Failed to ${action} task`);
    } finally {
      setTaskPendingAction(null);
    }
  };

  const handleStartTask = async (type: TaskType) => {
    try {
      setTaskPendingAction(`start-${type}`);
      await taskService.startTask(type);
      await fetchTask();
    } catch (err) {
      setError(err instanceof Error ? err.message : `Failed to start ${type} task`);
    } finally {
      setTaskPendingAction(null);
    }
  };

  const statItems = [
    {
      label: 'Total Rockets',
      value: stats.totalRockets,
      icon: Rocket,
      color: 'text-blue-600',
      path: '/rockets',
    },
    {
      label: 'Total Launches',
      value: stats.upcomingLaunches,
      icon: Calendar,
      color: 'text-green-600',
      path: '/launches',
    },
    {
      label: 'Launch Bases',
      value: stats.launchBases,
      icon: MapPin,
      color: 'text-orange-600',
      path: '/launch-bases',
    },
    {
      label: 'Agencies',
      value: stats.companies,
      icon: Building2,
      color: 'text-red-600',
      path: '/companies',
    },
    {
      label: 'News Articles',
      value: stats.newsArticles,
      icon: Newspaper,
      color: 'text-purple-600',
      path: '/news',
    },
  ];

  const ll2StatItems = [
    {
      label: 'LL2 Launchers',
      value: stats.ll2.launchers,
      icon: Rocket, // Reuse Rocket for now, or find another
      color: 'text-orange-600',
      path: '/ll2/launchers',
    },
    {
      label: 'LL2 Launches',
      value: stats.ll2.launches,
      icon: Rocket,
      color: 'text-blue-600',
      path: '/ll2/launches',
    },
    {
      label: 'LL2 Locations',
      value: stats.ll2.locations,
      icon: MapPin,
      color: 'text-green-600',
      path: '/ll2/locations',
    },
    {
      label: 'LL2 Agencies',
      value: stats.ll2.agencies,
      icon: Building2,
      color: 'text-red-600',
      path: '/ll2/agencies',
    },
    {
      label: 'LL2 Families',
      value: stats.ll2.launcher_families,
      icon: Layers,
      color: 'text-purple-600',
      path: '/ll2/families', // Verify path? Assume standard pattern for now
    },
    {
      label: 'LL2 Pads',
      value: stats.ll2.pads,
      icon: Target,
      color: 'text-pink-600',
      path: '/ll2/pads',
    },
  ];

  return (
    <div className="p-8">
      <h1 className="text-3xl font-bold mb-8">Dashboard</h1>

      {error && (
        <div className="mb-4 p-4 bg-red-100 text-red-800 rounded-md">
          {error}
        </div>
      )}

      {loading ? (
        <div className="text-center p-8">Loading dashboard statistics...</div>
      ) : (
        <>
          {showWelcome && (
            <Card className="mb-8">
              <CardHeader className="flex flex-row items-center justify-between">
                <CardTitle>Welcome to LaunchDate Admin</CardTitle>
                <Button
                  variant="ghost"
                  size="icon"
                  onClick={() => {
                    setShowWelcome(false);
                    localStorage.setItem('launchdate-admin-welcome-closed', 'true');
                  }}
                >
                  <X className="h-4 w-4" />
                </Button>
              </CardHeader>
              <CardContent>
                <p className="text-muted-foreground">
                  Manage your rocket launch data, news, companies, and launch bases from this admin panel.
                  Use the sidebar navigation to access different sections.
                </p>
              </CardContent>
            </Card>
          )}

          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-5 gap-6 mb-8">
            {statItems.map((stat) => {
              const Icon = stat.icon;
              return (
                <Card
                  key={stat.label}
                  className="cursor-pointer hover:shadow-lg transition-shadow"
                  onClick={() => navigate(stat.path)}
                >
                  <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                    <CardTitle className="text-sm font-medium">{stat.label}</CardTitle>
                    <Icon className={`h-4 w-4 ${stat.color}`} />
                  </CardHeader>
                  <CardContent>
                    <div className="text-2xl font-bold">{stat.value}</div>
                  </CardContent>
                </Card>
              );
            })}
          </div>

          <h2 className="text-2xl font-bold mb-4">LL2 Statistics</h2>
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-6 gap-6 mb-8">
            {ll2StatItems.map((stat) => {
              const Icon = stat.icon;
              return (
                <Card
                  key={stat.label}
                >
                  <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
                    <CardTitle className="text-sm font-medium">{stat.label}</CardTitle>
                    <Icon className={`h-4 w-4 ${stat.color}`} />
                  </CardHeader>
                  <CardContent>
                    <div className="text-2xl font-bold">{stat.value}</div>
                  </CardContent>
                </Card>
              );
            })}
          </div>


          <Card>
            <CardHeader className="flex flex-row items-center justify-between">
              <CardTitle>Task</CardTitle>
              <Button variant="outline" size="sm" onClick={fetchTask} disabled={loadingTask}>
                刷新
              </Button>
            </CardHeader>
            <CardContent className="space-y-4">
              {loadingTask ? (
                <div className="text-sm text-muted-foreground">Loading task...</div>
              ) : task ? (
                <div className="space-y-3">
                  <div className="grid gap-3 md:grid-cols-2">
                    <div>
                      <p className="text-xs text-muted-foreground">Type</p>
                      <p className="font-medium">{task.type}</p>
                    </div>
                    <div>
                      <p className="text-xs text-muted-foreground">Status</p>
                      <p className="font-medium">{task.status}</p>
                    </div>
                    <div>
                      <p className="text-xs text-muted-foreground">Started At</p>
                      <p className="font-medium">{new Date(task.started_at).toLocaleString()}</p>
                    </div>
                    <div>
                      <p className="text-xs text-muted-foreground">Updated At</p>
                      <p className="font-medium">{new Date(task.updated_at).toLocaleString()}</p>
                    </div>
                  </div>
                  {task.progress && (
                    <div>
                      <p className="text-xs text-muted-foreground mb-1">Progress</p>
                      <pre className="rounded-md bg-muted p-3 text-xs overflow-x-auto">
                        {JSON.stringify(task.progress, null, 2)}
                      </pre>
                    </div>
                  )}
                  {task.last_error && (
                    <div className="rounded-md border border-red-200 bg-red-50 p-3 text-sm text-red-700">
                      {task.last_error}
                    </div>
                  )}
                  <div className="flex flex-wrap gap-2">
                    <Button
                      variant="outline"
                      size="sm"
                      disabled={taskPendingAction !== null}
                      onClick={() => handleTaskAction('pause')}
                    >
                      Pause
                    </Button>
                    <Button
                      variant="outline"
                      size="sm"
                      disabled={taskPendingAction !== null}
                      onClick={() => handleTaskAction('resume')}
                    >
                      Resume
                    </Button>
                    <Button
                      variant="destructive"
                      size="sm"
                      disabled={taskPendingAction !== null}
                      onClick={() => handleTaskAction('cancel')}
                    >
                      Cancel
                    </Button>
                  </div>
                </div>
              ) : (
                <div className="text-sm text-muted-foreground">No running task.</div>
              )}

              <div>
                <p className="mb-2 text-sm font-medium">Start New Task</p>
                <div className="flex flex-wrap gap-2">
                  {TASK_TYPE_OPTIONS.map((item) => (
                    <Button
                      key={item.type}
                      variant="outline"
                      size="sm"
                      disabled={taskPendingAction !== null}
                      onClick={() => handleStartTask(item.type)}
                    >
                      {item.label}
                    </Button>
                  ))}
                </div>
              </div>
            </CardContent>
          </Card>
        </>
      )}
    </div>
  );
}
