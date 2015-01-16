%Testing WiFi tracking positioning based on Fingerprinting
%Binghao Li
%06/02/2014
%deterministic method, nearest neighbour
clear all;
close all;
load DB_143006022014;

meter_pixel = 0.0465; %0.0465/pixel
%ten test points in CE level4
true = [854 122; 789 122; 720 122; 655 122; 590 122; 522 122; 483 122; 386 122; 245 122; 162 122];  
%Android tablet
file_name = 'Testing_data_152006022014.csv';
%my iphone 4
%file_name = 'Testing_iphone_162711032014.csv';

fid = fopen(file_name, 'rt');
tmp = textscan(fid, '%d%d%d%s%d','HeaderLines', 1,'Delimiter',';');
fclose(fid);
ID = tmp{1};
Time_stamp = tmp{2};
infra_MAC = tmp{3};
MAC = tmp{4};
SS = tmp{5};

%SNAPSIMO5,6 and 8
AP_list = [5 6 8];
test_Fingerprint = [];
if(strcmp(file_name,'Testing_iphone_162711032014.csv'))

    b_temp = Time_stamp(2:end)-Time_stamp(1:end-1);
    Ia = find(b_temp>15);
    Ia = [0; Ia; length(Time_stamp)];    
    
    for i=1:length(Ia)-1
        SS_temp = SS(Ia(i)+1:Ia(i+1));
        tmp_fingerprint = [];
        for j=1:length(AP_list)
            I = find(infra_MAC(Ia(i)+1:Ia(i+1))==AP_list(j));
            if ~isempty(I)
                    tmp_fingerprint = [tmp_fingerprint mean(SS_temp(I)) length(SS_temp(I))];
                    %tmp_fingerprint = [tmp_fingerprint median(SS_temp(I)) length(SS_temp(I))];
            else
                    tmp_fingerprint = [tmp_fingerprint -100 0];
            end
        end
        %format: x(pixel, true) y(pixel, true) mean_SS no_scan(SNAPSIMO05) mean_SS no_scan(SNAPSIMO06) mean_SS no_scan(SNAPSIMO08)
        test_Fingerprint = [test_Fingerprint; tmp_fingerprint];
    end

    
else
    scan_time_file_name = 'test_scan_time_152006022014.txt';
    fid_scan_time = fopen(scan_time_file_name, 'rt');
    tline = fgetl(fid_scan_time);
    ended = 0;
    scan_time = [];
   

    while ischar(tline)
        while ~ended
            %disp(tline)
            if (strcmp(tline(1:5),'START') == 1)
                temp = regexp(tline,'@','split');
                startPoint = sscanf(temp{2}, '(%d,%d)');
                startTime = str2double(temp{3});
            elseif (strcmp(tline(1:3),'END') == 1)
                temp = regexp(tline,'@','split');
                endPoint = sscanf(temp{2}, '(%d,%d)');
                endTime = str2double(temp{3});
                ended = 1;
            else
                temp = regexp(tline,'@','split');
                scan_time = [scan_time; str2double(temp{2})];
            end
            tline = fgetl(fid);
        end

        ended = 0;

        %find the SSs recorded by APs, there is short delay (not delay but the synchronization issue), add 1 or 2 seconds to
        %make sure the SSs can be found
        temp_I = find(Time_stamp>=endTime-1 & Time_stamp<=endTime+2);
        SS_temp = SS(temp_I);
        tmp_fingerprint = [];
        for j=1:length(AP_list)
            I = find(infra_MAC(temp_I)==AP_list(j));
            if ~isempty(I)
                tmp_fingerprint = [tmp_fingerprint mean(SS_temp(I)) length(SS_temp(I))];
                %tmp_fingerprint = [tmp_fingerprint median(SS_temp(I)) length(SS_temp(I))];
            else
                tmp_fingerprint = [tmp_fingerprint -100 0];
            end
        end
        %format: x(pixel, true) y(pixel, true) mean_SS no_scan(SNAPSIMO05) mean_SS no_scan(SNAPSIMO06) mean_SS no_scan(SNAPSIMO08) 
        test_Fingerprint = [test_Fingerprint; tmp_fingerprint];
    end
end

test_Fingerprint = [true test_Fingerprint];

%% Positioning, nearest neighbour

tmp1 = double(Fingerprint(:,3:2:end));
estimate = [];
for i=1:length(true(:,1))
    tmp2 = ones(length(Fingerprint(:,1)),1)*double(test_Fingerprint(i,3:2:end));
    Euc_dis = sum((tmp2-tmp1).^2,2);
    [C,I] = min(Euc_dis);
    %I
    estimate = [estimate; Fingerprint(I,1:2)];
end
error = abs(double(estimate) - true)*meter_pixel


%% Positioning, trilateration, using an old model, see my thesis pp50
%{
%r_n=r_0*10^((RSS_0-RSS)/10*alfa)
%y=ax+b, b=RSS0, a=-10*alfa, x=log10 r_n/r_0 
%Linear fit, r0=2.25, RSS0=-51.14, 10*alfa=11.098

AP05 = [885 320]*meter_pixel;
AP06 = [511 70]*meter_pixel;
AP08 = [154 191]*meter_pixel;

r0 = 2.25;
RSS0 = -51.14;
alfa10 = 11.098;

estimate=[];
for i=1:length(test_Fingerprint(:,1))
    RSS05 = test_Fingerprint(i,3);
    %r05 = r0*10^((RSS0-RSS05)/alfa10);
    r05 = 0.000198*abs(RSS05)^3-0.025*RSS05^2+1.14*abs(RSS05)-14.8;
    RSS06 = test_Fingerprint(i,5);
    %r06 = r0*10^((RSS0-RSS06)/alfa10);
    r06 = 0.000198*abs(RSS06)^3-0.025*RSS06^2+1.14*abs(RSS06)-14.8;
    RSS08 = test_Fingerprint(i,7);
    %r08 = r0*10^((RSS0-RSS08)/alfa10);
    r08 = 0.000198*abs(RSS08)^3-0.025*RSS08^2+1.14*abs(RSS08)-14.8;
    
    [r05 r06 r08]
    
    circle(AP05(1),AP05(2),r05);
    hold on;
    circle(AP06(1),AP06(2),r06);
    circle(AP08(1),AP08(2),r08);
    hold off;
    xy = geometric_trilateration(AP05, AP06, AP08, r05, r06, r08);
    estimate = [estimate; xy];
end
true = true*meter_pixel;
%}

a=1;
%[xout,yout] = circcirc(x1,y1,r1,x2,y2,r2)





