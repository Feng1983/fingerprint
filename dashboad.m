%Analyse data, apply fingerprinting to locate the users,
%generate several plots to visualize the results
%Binghao Li 
%27/03/2014
clear all;
close all;

WeakSignal = -100;
Hour = 3600;
Day = 86400;
Week = 604800;

MyiPhone4='50EAD68D0A17'; %dec: 88969552136727
Samsung_nexus='B407F947D4AB'; %dec: 197946340005035
SamsungGalaxyNexusS='A00bba2d85cd'; %dec: 175972228629965
AndroidTablet='08606E388A2B'; %dec: 9210259081771
MyUbuntu='001B772CB037'; %dec: 117963534391
YincaiiPhone='C063944dd860'; %dec: 211533922424928
SIMOiphone='18e7f4e8c336'; %dec: 27384525407030

file_name = 'CELevel4_000024032014_230006042014local.csv'; %note, 05/04 daylight saving was ended, that is why 23:00 was used as end time.
%file_name = 'CELevel4_000024032014_000025032014.csv';
%file_name = 'AndroidTablet120001042014.csv';
%file_name = 'SamsungNexusS06001042014.csv';
fid = fopen(file_name, 'rt');
tmp = textscan(fid, '%d%d%d%d64%d%d','HeaderLines', 1,'Delimiter',';');
fclose(fid);
ID = tmp{1};
Time_stamp = tmp{2};
infra_MAC = tmp{3};
MAC = tmp{4};
temp = dec2hex(MAC);
MAC_string = cellstr(temp);
SS = tmp{5};
Frequency = tmp{6};

%% Apply fingerprinting to localize the mobile devices
time_ref = Time_stamp(1);
time_end = Time_stamp(end);

tmp_data = [];
tmp_MAC = [];
PP_test = [];

for i=1:1:time_end-time_ref
    if mod(i,10000)==0
       i 
    end
    %note, use 2 seconds as a group, there are over laps
    I = find(Time_stamp >= (time_ref + i - 1) & Time_stamp <= (time_ref + i));
    nID = ID(I);
    nTime_stamp = Time_stamp(I);
    ninfra_MAC = infra_MAC(I);
    nMAC = MAC(I);
    nMAC_string = MAC_string(I,:);
    nSS = SS(I);
    nFrequency = Frequency(I);
    %[uni_MAC,junk,ind] = unique(nMAC_string);

    I5 = find(ninfra_MAC==5);
    I6 = find(ninfra_MAC==6);
    I8 = find(ninfra_MAC==8);

    SNAPSIMO05 = [];
    SNAPSIMO06 = [];
    SNAPSIMO08 = [];
    SNAPSIMO05_MAC = [];
    SNAPSIMO06_MAC = [];
    SNAPSIMO08_MAC = [];
    
    if(~isempty(I5))
        SNAPSIMO05 = double([nTime_stamp(I5) nSS(I5) nFrequency(I5)]);
        SNAPSIMO05_MAC = nMAC(I5);
        %SNAPSIMO05_MAC_string = nMAC_string(I5);
    end
    if(~isempty(I6))
        SNAPSIMO06 = double([nTime_stamp(I6) nSS(I6) nFrequency(I6)]);
        SNAPSIMO06_MAC = nMAC(I6);
        %SNAPSIMO06_MAC_string = nMAC_string(I6);
    end
    if(~isempty(I8)) 
        SNAPSIMO08 = double([nTime_stamp(I8) nSS(I8) nFrequency(I8)]);
        SNAPSIMO08_MAC = nMAC(I8);
        %SNAPSIMO08_MAC_string = nMAC_string(I8);
    end
    %find member in cell array
    %[truefalse, index] = ismember('string', cell_array);
    %[rn,cn]=find(strcmp(cell_array,'sring'));
    
     
    %test a way to find the MACs in CE level4:
    %1. at least two APs detected signals
    %2. only one AP detected strong signal
    %remove those MACs
    %1.three APs detect very weak signals
    %2.two APs detect very weak signals
    uni_MAC = unique([SNAPSIMO05_MAC; SNAPSIMO06_MAC; SNAPSIMO08_MAC]);
    
    
    if(~isempty(uni_MAC))
        for j=1:length(uni_MAC)
            %find those MAC detected by at least two APs
            tmp_I = find(SNAPSIMO05_MAC==uni_MAC(j));
            if(~isempty(tmp_I))
                %the MAC may detected twice in 2 second, get the average value
                ss05 = [sum(SNAPSIMO05(tmp_I,2).*SNAPSIMO05(tmp_I,3))/sum(SNAPSIMO05(tmp_I,3)) sum(SNAPSIMO05(tmp_I,3))];
            else
                ss05 = [WeakSignal 0];//??-100??
            end

            tmp_I = find(SNAPSIMO06_MAC==uni_MAC(j));
            if(~isempty(tmp_I))
                ss06 = [sum(SNAPSIMO06(tmp_I,2).*SNAPSIMO06(tmp_I,3))/sum(SNAPSIMO06(tmp_I,3)) sum(SNAPSIMO06(tmp_I,3))];
            else
                ss06 = [WeakSignal 0];
            end

            tmp_I = find(SNAPSIMO08_MAC==uni_MAC(j));
            if(~isempty(tmp_I))
                ss08 = [sum(SNAPSIMO08(tmp_I,2).*SNAPSIMO08(tmp_I,3))/sum(SNAPSIMO08(tmp_I,3)) sum(SNAPSIMO08(tmp_I,3))];
            else
                ss08 = [WeakSignal 0];
            end
            
            %find those signals from level 4 only, details for these
            %criteria can be found in Test Result.doc,section 'Different
            %floors tests'.
            %There is no AP at levels other than level 4, it is not easy to
            %tell if a device is in level 4 or not, if there are APs on
            %other levels, this task will be very simple.
            
            %find at least one less very strong signal (>-65)
            if ~isempty(find([ss05(1) ss06(1) ss08(1)]>-65, 1))
                tmp_MAC = [tmp_MAC; uni_MAC(j)];
                tmp_data = [tmp_data; double(time_ref+i-1) ss05 ss06 ss08]; 
            %two_strong
            elseif length(find([ss05(1) ss06(1) ss08(1)]>-70))==2
                tmp_MAC = [tmp_MAC; uni_MAC(j)];
                tmp_data = [tmp_data; double(time_ref+i-1) ss05 ss06 ss08]; 
            %one_strong_one_med (<=-70 >=-80) & one_strong_two_med    
            elseif (length(find([ss05(1) ss06(1) ss08(1)]>-70))==1 && length(find([ss05(1) ss06(1) ss08(1)]>=-80 & [ss05(1) ss06(1) ss08(1)]<=-70))>=1)
                tmp_MAC = [tmp_MAC; uni_MAC(j)];
                tmp_data = [tmp_data; double(time_ref+i-1) ss05 ss06 ss08]; 
            %three_med or two_med
            elseif length(find([ss05(1) ss06(1) ss08(1)]<=-70 & [ss05(1) ss06(1) ss08(1)]>=-80))>=2
                tmp_MAC = [tmp_MAC; uni_MAC(j)];
                tmp_data = [tmp_data; double(time_ref+i-1) ss05 ss06 ss08];                      
            end
         end
    end
   
end

%save CELevel4_000024032014_230006042014local PP_test tmp_MAC tmp_data;
%load CELevel4_000024032014_230006042014local;
%Applying fingerprinting, nearest neighbour
load CE_level4_FPDB.mat;
tmp1 = double(Fingerprint(:,3:2:end));
estimate = [];
tmp_C = [];
for i=1:length(tmp_MAC)
    %the positioning algorithm should be refined. If there are only one
    %strong signla detected, it might be better to apply another algorithm
    %rather than NN!!!!
    if mod(i,10000)==0
       i 
    end
    tmp2 = ones(length(Fingerprint(:,1)),1)*tmp_data(i,2:2:6);
    Euc_dis = sum((tmp2-tmp1).^2,2);
    [C,I] = min(Euc_dis);
    tmp_C = [tmp_C; C];
    estimate = [estimate; Fingerprint(I,1:2)];
    %KNN?
    %[C,I] = sort(Euc_dis);
    %estimate = [estimate; mean(Fingerprint(I(1:3),1:2),1)];
end


%{
%find unique MAC first
%then find the frequency of each MAC
%and sorted from low frequency to high requency
[tMAC,junk,ind] = unique(tmp_MAC);
freq_MAC = histc(ind,1:numel(tMAC));
[sorted_freq_MAC, Ind_freq_MAC] = sort(freq_MAC);

%to find the most frequent appeared MAC and get the index
MAC2show = tMAC(Ind_freq_MAC(end));
tmp_I = find(tmp_MAC == MAC2show);
%}

%show the floor plan
figure;
I = imread('CE_level4.png'); 
imshow(I);
hold on;
%{
%display the position of chosen MAC in time series
for k=1:length(tmp_I)
    plot(estimate(tmp_I(k),1),estimate(tmp_I(k),2),'*');
    pause(0.5);
    plot(estimate(tmp_I(k),1),estimate(tmp_I(k),2),'w*');
    pause(0.5);
end
%}

%the same RP may have slight different coordiantes (because of the way to
%creat the database, for instance, from A to B, one RP is created, from B
%to A, another RP is created, these two should be the same RP,
%however, the coordinate may slightly diffrent (in double), if ignor the
%deciaml points, they are the same
int_estimate = int32(estimate);

[ua,ja,ia]=unique(int_estimate,'rows');
num_person = [];
for k=1:1:length(ua)
   iaI=find(ia==k);
   p_MAC = tmp_MAC(iaI); %find the MACs located at the same location
   num_person = [num_person length(unique(p_MAC))/5]; %number of unique MAC gives the number of person at this location for a specific period of time, 5-10 minutes? 
end


for i=1:length(ua)
    plot(ua(i,1),ua(i,2),'*','MarkerSize', num_person(i));
end


%Generate visitor number
%date = datestr(time_ref/86400 + datenum(1970,1,1))

%% Process the data NOT based on position - number of visitors etc
if 1==1
%Visitors per hour or per day or per week
%step = Hour/Day/Week
step = Hour;
no_visitors = [];

%this is used to find the MAC can be seen all the time. Can be used to find
%the staff working in a shop, devices (such as laptop) other than a mobile
%phone etc.

MAC_can_be_seen_allperiod = [];

for i=time_ref:step:time_end-1
    
    tmp_I = find(i<=tmp_data(:,1)& tmp_data(:,1)<i+step);
    if i==time_ref
        MAC_can_be_seen_allperiod = unique(tmp_MAC(tmp_I));
    else
        MAC_can_be_seen_allperiod = intersect(unique(tmp_MAC(tmp_I)), MAC_can_be_seen_allperiod);
    end
    no_visitors = [no_visitors; length(unique(tmp_MAC(tmp_I)))]; %for a shop, this number - MAC_can_be_seen_allperiod can be the real number of visitors
    
end
%Visitors per day
figure;
bar(no_visitors);

end

a=1;
%Visitors per week etc.

%% Check the traffic using MAC only - MACs can be seen by all three APs during a period (say 1 hour)
% or MACs can be seen by at least two APs during a period
% This part can be ignored, used for test only
if 1==2
time_ref = Time_stamp(1);
traffic = [];
MAC_Cn_by3stations_all=[];
step = 1;
t_step = step*3600;
for i=0:step:23
    
    I = find(Time_stamp>=(time_ref+i*3600) & Time_stamp<(time_ref+i*3600+t_step));
    nID = ID(I);
    nTime_stamp = Time_stamp(I);
    ninfra_MAC = infra_MAC(I);
    nMAC = MAC(I);
    nSS = SS(I);

    I5 = find(ninfra_MAC==5);
    I6 = find(ninfra_MAC==6);
    I8 = find(ninfra_MAC==8);

    SNAPSIMO05_Time_stamp = nTime_stamp(I5);
    SNAPSIMO05_MAC = nMAC(I5);
    SNAPSIMO05_SS = nSS(I5);

    SNAPSIMO06_Time_stamp = nTime_stamp(I6);
    SNAPSIMO06_MAC = nMAC(I6);
    SNAPSIMO06_SS = nSS(I6);

    SNAPSIMO08_Time_stamp = nTime_stamp(I8);
    SNAPSIMO08_MAC = nMAC(I8);
    SNAPSIMO08_SS = nSS(I8);

    temp1 = intersect(SNAPSIMO05_MAC, SNAPSIMO06_MAC);
    temp2 = intersect(SNAPSIMO05_MAC, SNAPSIMO08_MAC);
    temp3 = intersect(SNAPSIMO06_MAC, SNAPSIMO08_MAC);
    [uni_MAC,junk,ind] = unique([temp1; temp2; temp3]);
    
    MAC_Cn_by3stations = intersect(temp1, SNAPSIMO08_MAC);
    MAC_Cn_by3stations_all = [MAC_Cn_by3stations_all;  MAC_Cn_by3stations];
    %traffic = [traffic length(MAC_Cn_by3stations)];
    traffic = [traffic length(uni_MAC)];
end

figure;
bar(traffic);

end

a=1;

