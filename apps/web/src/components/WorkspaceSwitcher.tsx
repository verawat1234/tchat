import React from 'react';
import { Button } from './ui/button';
import { Badge } from './ui/badge';
import { Avatar, AvatarFallback, AvatarImage } from './ui/avatar';
import { DropdownMenu, DropdownMenuContent, DropdownMenuItem, DropdownMenuLabel, DropdownMenuSeparator, DropdownMenuTrigger } from './ui/dropdown-menu';
import { 
  Building2, Store, Users, Crown, ArrowDown, Plus, Briefcase,
  TrendingUp, DollarSign, MessageCircle, ChevronDown
} from 'lucide-react';
import { toast } from "sonner";

interface Workspace {
  id: string;
  name: string;
  type: 'shop' | 'company';
  role: string;
  avatar: string;
  customerCount: number;
  revenue?: string;
}

interface WorkspaceSwitcherProps {
  selectedWorkspace: string;
  userWorkspaces: Workspace[];
  onWorkspaceChange: (workspaceId: string) => void;
  variant?: 'default' | 'compact' | 'prominent';
  showMetrics?: boolean;
}

export function WorkspaceSwitcher({ 
  selectedWorkspace, 
  userWorkspaces, 
  onWorkspaceChange,
  variant = 'default',
  showMetrics = true
}: WorkspaceSwitcherProps) {
  const currentWorkspace = userWorkspaces.find(w => w.id === selectedWorkspace);

  if (!currentWorkspace) {
    return null;
  }

  const getWorkspaceIcon = (type: string, size = 'w-4 h-4') => {
    switch (type) {
      case 'shop':
        return <Store className={`${size} text-chart-2`} />;
      case 'company':
        return <Building2 className={`${size} text-chart-3`} />;
      default:
        return <Briefcase className={`${size} text-muted-foreground`} />;
    }
  };

  const getRoleIcon = (role: string) => {
    if (role.toLowerCase().includes('owner')) {
      return <Crown className="w-3 h-3 text-yellow-500" />;
    }
    return null;
  };

  const handleAddWorkspace = () => {
    toast.success('Opening workspace creation wizard...');
  };

  // Prominent variant for top of Work tab
  if (variant === 'prominent') {
    return (
      <div className="bg-gradient-to-r from-chart-1/10 to-chart-2/10 border-b border-border p-4">
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <Button variant="ghost" className="w-full h-auto p-4 hover:bg-accent/50 transition-all duration-200 hover:scale-[1.02]">
              <div className="flex items-center gap-4 w-full">
                <div className="relative">
                  <Avatar className="w-16 h-16 border-2 border-background shadow-lg">
                    <AvatarImage src={currentWorkspace.avatar} />
                    <AvatarFallback className="text-lg">
                      {getWorkspaceIcon(currentWorkspace.type, 'w-8 h-8')}
                    </AvatarFallback>
                  </Avatar>
                  <div className="absolute -bottom-1 -right-1 bg-green-500 w-5 h-5 rounded-full border-2 border-background flex items-center justify-center">
                    {getRoleIcon(currentWorkspace.role) || <div className="w-2 h-2 bg-white rounded-full" />}
                  </div>
                </div>
                
                <div className="flex-1 text-left min-w-0">
                  <div className="flex items-center gap-2 mb-1">
                    <h3 className="font-semibold text-lg truncate">{currentWorkspace.name}</h3>
                    <ChevronDown className="w-5 h-5 text-muted-foreground" />
                  </div>
                  <div className="flex items-center gap-2 mb-2">
                    <Badge variant="secondary" className="text-xs">
                      {currentWorkspace.role}
                    </Badge>
                    <Badge variant="outline" className="text-xs">
                      {currentWorkspace.type === 'shop' ? 'Restaurant' : 'Company'}
                    </Badge>
                  </div>
                  
                  {showMetrics && (
                    <div className="flex items-center gap-4 text-sm text-muted-foreground">
                      <div className="flex items-center gap-1">
                        <Users className="w-4 h-4" />
                        <span>{currentWorkspace.customerCount} customers</span>
                      </div>
                      {currentWorkspace.revenue && (
                        <div className="flex items-center gap-1">
                          <TrendingUp className="w-4 h-4 text-green-500" />
                          <span className="text-green-600 font-medium">{currentWorkspace.revenue}</span>
                        </div>
                      )}
                    </div>
                  )}
                </div>
              </div>
            </Button>
          </DropdownMenuTrigger>
          
          <DropdownMenuContent align="start" className="w-80">
            <DropdownMenuLabel className="flex items-center gap-2 py-3">
              <Briefcase className="w-5 h-5" />
              Switch Workspace
            </DropdownMenuLabel>
            <DropdownMenuSeparator />
            
            {userWorkspaces.map((workspace) => (
              <DropdownMenuItem
                key={workspace.id}
                onClick={() => onWorkspaceChange(workspace.id)}
                className={`p-4 cursor-pointer ${
                  workspace.id === selectedWorkspace ? 'bg-accent' : ''
                }`}
              >
                <div className="flex items-center gap-3 w-full">
                  <div className="relative">
                    <Avatar className="w-12 h-12">
                      <AvatarImage src={workspace.avatar} />
                      <AvatarFallback>
                        {getWorkspaceIcon(workspace.type, 'w-6 h-6')}
                      </AvatarFallback>
                    </Avatar>
                    {workspace.id === selectedWorkspace && (
                      <div className="absolute -bottom-1 -right-1 w-4 h-4 bg-primary rounded-full border-2 border-background"></div>
                    )}
                  </div>
                  
                  <div className="flex-1 min-w-0">
                    <div className="flex items-center gap-2 mb-1">
                      <span className="font-medium truncate">{workspace.name}</span>
                      {getRoleIcon(workspace.role)}
                    </div>
                    <div className="flex items-center gap-2 mb-1">
                      <Badge variant="secondary" className="text-xs">{workspace.role}</Badge>
                    </div>
                    <div className="flex items-center gap-3 text-xs text-muted-foreground">
                      <span>{workspace.customerCount} customers</span>
                      {workspace.revenue && (
                        <span className="text-green-600 font-medium">{workspace.revenue}</span>
                      )}
                    </div>
                  </div>
                </div>
              </DropdownMenuItem>
            ))}
            
            <DropdownMenuSeparator />
            <DropdownMenuItem onClick={handleAddWorkspace} className="p-3">
              <Plus className="w-4 h-4 mr-2" />
              Add Workspace
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      </div>
    );
  }

  // Compact variant for header
  if (variant === 'compact') {
    return (
      <DropdownMenu>
        <DropdownMenuTrigger asChild>
          <Button variant="ghost" size="sm" className="h-8 gap-2">
            <Avatar className="w-6 h-6">
              <AvatarImage src={currentWorkspace.avatar} />
              <AvatarFallback className="text-xs">
                {getWorkspaceIcon(currentWorkspace.type, 'w-3 h-3')}
              </AvatarFallback>
            </Avatar>
            <span className="font-medium text-sm truncate max-w-24">
              {currentWorkspace.name}
            </span>
            <ArrowDown className="w-3 h-3" />
          </Button>
        </DropdownMenuTrigger>
        
        <DropdownMenuContent align="end" className="w-64">
          <DropdownMenuLabel>Switch Workspace</DropdownMenuLabel>
          <DropdownMenuSeparator />
          
          {userWorkspaces.map((workspace) => (
            <DropdownMenuItem
              key={workspace.id}
              onClick={() => onWorkspaceChange(workspace.id)}
              className={workspace.id === selectedWorkspace ? 'bg-accent' : ''}
            >
              <Avatar className="w-8 h-8 mr-2">
                <AvatarImage src={workspace.avatar} />
                <AvatarFallback>
                  {getWorkspaceIcon(workspace.type, 'w-4 h-4')}
                </AvatarFallback>
              </Avatar>
              <div className="flex-1">
                <p className="font-medium">{workspace.name}</p>
                <p className="text-xs text-muted-foreground">{workspace.role}</p>
              </div>
            </DropdownMenuItem>
          ))}
        </DropdownMenuContent>
      </DropdownMenu>
    );
  }

  // Default variant
  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <Button variant="outline" className="h-auto p-3 justify-start gap-3 hover:bg-accent/50">
          <Avatar className="w-10 h-10">
            <AvatarImage src={currentWorkspace.avatar} />
            <AvatarFallback>
              {getWorkspaceIcon(currentWorkspace.type, 'w-5 h-5')}
            </AvatarFallback>
          </Avatar>
          
          <div className="flex-1 text-left min-w-0">
            <div className="flex items-center gap-2 mb-1">
              <span className="font-medium truncate">{currentWorkspace.name}</span>
              {getRoleIcon(currentWorkspace.role)}
            </div>
            <div className="flex items-center gap-2">
              <Badge variant="secondary" className="text-xs">{currentWorkspace.role}</Badge>
              {showMetrics && currentWorkspace.revenue && (
                <span className="text-xs text-green-600 font-medium">{currentWorkspace.revenue}</span>
              )}
            </div>
          </div>
          
          <ArrowDown className="w-4 h-4 text-muted-foreground" />
        </Button>
      </DropdownMenuTrigger>
      
      <DropdownMenuContent align="start" className="w-72">
        <DropdownMenuLabel className="flex items-center gap-2">
          <Briefcase className="w-4 h-4" />
          Your Workspaces
        </DropdownMenuLabel>
        <DropdownMenuSeparator />
        
        {userWorkspaces.map((workspace) => (
          <DropdownMenuItem
            key={workspace.id}
            onClick={() => onWorkspaceChange(workspace.id)}
            className={`p-3 cursor-pointer ${
              workspace.id === selectedWorkspace ? 'bg-accent' : ''
            }`}
          >
            <div className="flex items-center gap-3 w-full">
              <Avatar className="w-10 h-10">
                <AvatarImage src={workspace.avatar} />
                <AvatarFallback>
                  {getWorkspaceIcon(workspace.type, 'w-5 h-5')}
                </AvatarFallback>
              </Avatar>
              
              <div className="flex-1 min-w-0">
                <div className="flex items-center gap-2 mb-1">
                  <span className="font-medium truncate">{workspace.name}</span>
                  {getRoleIcon(workspace.role)}
                  {workspace.id === selectedWorkspace && (
                    <Badge className="text-xs">Active</Badge>
                  )}
                </div>
                <div className="flex items-center gap-2 text-xs text-muted-foreground">
                  <span>{workspace.role}</span>
                  <span>•</span>
                  <span>{workspace.customerCount} customers</span>
                  {workspace.revenue && (
                    <>
                      <span>•</span>
                      <span className="text-green-600 font-medium">{workspace.revenue}</span>
                    </>
                  )}
                </div>
              </div>
            </div>
          </DropdownMenuItem>
        ))}
        
        <DropdownMenuSeparator />
        <DropdownMenuItem onClick={handleAddWorkspace}>
          <Plus className="w-4 h-4 mr-2" />
          Create New Workspace
        </DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>
  );
}