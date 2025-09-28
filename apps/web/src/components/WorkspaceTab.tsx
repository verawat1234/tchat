/**
 * Temporary Workspace Tab - Simple version for testing gateway integration
 */

import React from 'react';

const WorkspaceTab: React.FC = () => {
  return (
    <div className="h-full flex items-center justify-center bg-background">
      <div className="text-center">
        <h2 className="text-2xl font-bold mb-4">Workspace</h2>
        <p className="text-muted-foreground">
          Gateway integration testing in progress...
        </p>
      </div>
    </div>
  );
};

export default WorkspaceTab;